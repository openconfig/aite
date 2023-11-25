// Copyright 2023 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package srv implements a server for the Aite service. Aite manipulates the
// state and performance of interfaces within a container's network namespace.
package srv

import (
	"context"
	"fmt"
	"math"
	"net"
	"time"

	"golang.org/x/sys/unix"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"

	"github.com/florianl/go-tc"
	"github.com/florianl/go-tc/core"
	"github.com/openconfig/magna/intf"

	apb "github.com/openconfig/aite/proto/aite"
)

// S is the wrapper for the Aite service implementation.
type S struct {
	tc *tc.Tc

	*apb.UnimplementedAiteServer
}

// New returns a new Aite server.
func New() (*S, error) {
	tconn, err := tc.Open(&tc.Config{})
	if err != nil {
		return nil, fmt.Errorf("cannot open Tc connection, %w", err)
	}

	return &S{tc: tconn}, nil
}

// Stop stops the Aite server, cleaning up internal state.
func (s *S) Stop() error {
	if err := s.tc.Close(); err != nil {
		return fmt.Errorf("cannot close Tc connection, %v", err)
	}
	return nil
}

// SetInterfaceState implements the InterfaceState RPC for the Aite service. It
// manipulates parameters of the interface including impairments.
func (s *S) SetInterface(ctx context.Context, req *apb.SetInterfaceRequest) (*apb.SetInterfaceResponse, error) {
	if req.Name == "" || !intf.ValidInterface(req.Name) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid interface name specified, %s", req.Name)
	}

	if req.GetParams() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "params is a required argument")
	}

	params := req.GetParams()
	if params.State == apb.InterfaceState_IS_UNSPECIFIED {
		return nil, status.Errorf(codes.InvalidArgument, "interface state must be specified")
	}

	var iState intf.IntState
	switch params.State {
	case apb.InterfaceState_IS_UP:
		iState = intf.InterfaceUp
	case apb.InterfaceState_IS_ADMIN_DOWN:
		iState = intf.InterfaceDown
	default:
		return nil, status.Errorf(codes.InvalidArgument, "invalid interface state %s specified", params.State)
	}

	if params.LossPct >= 100 {
		return nil, status.Errorf(codes.InvalidArgument, "loss percentage must be 0 < loss <= 100, got: %d", params.LossPct)
	}

	if err := s.applyInterfaceState(ctx, req.Name, iState, params.LossPct, params.LatencyMsec); err != nil {
		return nil, status.Errorf(codes.Internal, "cannot set interface state, %v", err)
	}

	return &apb.SetInterfaceResponse{
		Name: req.Name,
		Params: &apb.InterfaceStateParams{
			State:       req.GetParams().GetState(),
			LatencyMsec: req.GetParams().GetLatencyMsec(),
			LossPct:     req.GetParams().GetLossPct(),
		},
	}, nil
}

// applyInterfaceState applies the state changes to the interface with the specified name. The
// state indicates any change in administrative or operational status. lossPct indicates a percentage
// latency that should be applied, and latencyMsec is an additional latency that should be applied
// to packets traversing the interface.
func (s *S) applyInterfaceState(ctx context.Context, name string, state intf.IntState, lossPct, latencyMsec uint32) error {
	if err := intf.InterfaceState(name, state); err != nil {
		return status.Errorf(codes.Internal, "cannot set interface state, %v", err)
	}

	if err := s.impairInterface(ctx, name, lossPct, latencyMsec); err != nil {
		return err
	}

	return nil
}

const (
	// MaxUint32 is the maximum uint32 value. In netem a number of values are
	// expressed by indicating the proportion of max uint32.
	MaxUint32 uint32 = 0xFFFFFFFF
)

// impairInterface applies the specified impairments to the interface with the specified name. lossPct
// indicates a percentage packet loss to be applied, and latencyMsec an additional artificial latency
// to be applied. The function will set the underlying kernel parameters regardless of their current
// state.
func (s *S) impairInterface(ctx context.Context, name string, lossPct, latencyMsec uint32) error {
	intID, err := net.InterfaceByName(name)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "cannot find interface %s", name)
	}

	qopt := tc.NetemQopt{
		// Maximum number of packets that should be in the buffer,
		// the default value of 0 will stop packets flowing.
		Limit: 1000,
	}

	// latency is set in µsec, and is fed to the kernel as CPU
	// ticks, not msec. Thus, we need to use the tc package's
	// helper to convert to ticks.
	qopt.Latency = core.Time2Tick(uint32(latencyMsec * 1000))
	klog.Infof("setting device %s latency to %d msec", name, latencyMsec)

	p := uint32(math.Round(float64(MaxUint32) * (float64(lossPct) / 100.0)))
	qopt.Loss = p
	klog.Infof("setting device %s loss to %d%% (val: %d)", name, lossPct, p)

	qdisc := tc.Object{
		Msg: tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(intID.Index),
			Handle:  core.BuildHandle(0x1, 0x0),
			Parent:  tc.HandleRoot,
			Info:    0,
		},
		Attribute: tc.Attribute{
			Kind: "netem",
			Netem: &tc.Netem{
				Qopt: qopt,
			},
		},
	}

	// We should not ever block on the qdisc call below, but to ensure that we have a
	// reasonable belt and braces approach here, we use the parent context to ensure
	// that we cancel the context if we do.
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	klog.Infof("calling qdisc replace")
	if err := s.tc.Qdisc().Replace(&qdisc); err != nil {
		return status.Errorf(codes.Internal, "cannot apply impairment to interface, %v", err)
	}
	klog.Infof("returned from qdisc replace")

	return nil
}
