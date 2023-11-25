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
	"k8s.io/klog/v2"

	"github.com/florianl/go-tc"
	"github.com/florianl/go-tc/core"
	"github.com/openconfig/magna/intf"

	apb "github.com/robshakir/aite/proto/aite"
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

// InterfaceState implements the InterfaceState method of the Aite service. It
// shuts down the interface specified in the request.
func (s *S) InterfaceState(_ context.Context, req *apb.InterfaceStateRequest) (*apb.InterfaceStateResponse, error) {
	state := intf.InterfaceUp
	if req.Shutdown {
		state = intf.InterfaceDown
	}
	if err := intf.InterfaceState(req.Name, state); err != nil {
		return nil, status.Errorf(codes.Internal, "cannot set interface state, %v", err)
	}

	return &apb.InterfaceStateResponse{}, nil
}

const (
	// MaxUint32 is the maximum uint32 value. In netem a number of values are
	// expressed by indicating the proportion of max uint32.
	MaxUint32 uint32 = 0xFFFFFFFF
)

// ImpairInterface implements the ImpairInterface method mof the Aite service. It
// adds the specified impairment to the interface specified in the request.
func (s *S) ImpairInterface(ctx context.Context, req *apb.ImpairInterfaceRequest) (*apb.ImpairInterfaceResponse, error) {
	intID, err := net.InterfaceByName(req.Name)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "cannot find interface %s", req.Name)
	}

	qopt := tc.NetemQopt{
		// Maximum number of packets that should be in the buffer,
		// the default value of 0 will stop packets flowing.
		Limit: 1000,
	}

	if req.GetLatencyMsec() != 0 {
		// latency is set in µsec, and is fed to the kernel as CPU
		// ticks, not msec. Thus, we need to use the tc package's
		// helper to convert to ticks.
		qopt.Latency = core.Time2Tick(uint32(req.LatencyMsec * 1000))
		klog.Infof("setting device %s latency to %d msec", req.Name, req.LatencyMsec)
	}

	if pct := req.GetLossPct(); pct != 0 {
		p := uint32(math.Round(float64(MaxUint32) * (float64(pct) / 100.0)))
		qopt.Loss = p
		klog.Infof("setting device %s loss to %d%% (val: %d)", req.Name, pct, p)
	}

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

	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	klog.Infof("calling qdisc replace")
	if err := s.tc.Qdisc().Replace(&qdisc); err != nil {
		return nil, status.Errorf(codes.Internal, "cannot apply impairment to interface, %v", err)
	}
	klog.Infof("returned from qdisc replace")

	return &apb.ImpairInterfaceResponse{}, nil
}
