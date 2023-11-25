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

// Binary aite implements a container that can manipulate the state and performance
// of interfaces within its local network namespace, controlled by a gRPC API. This
// provides a mechanism for an external caller to introduce impairments to a network
// simulation (such as one started by KNE).
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/openconfig/aite/srv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"k8s.io/klog/v2"

	apb "github.com/openconfig/aite/proto/aite"
)

var (
	port = flag.Uint("port", 60061, "port for the aite service to listen on")
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	serv := grpc.NewServer()
	as, err := srv.New()
	if err != nil {
		klog.Exitf("cannot create Aite server, %v", err)
	}
	defer func() {
		if err := as.Stop(); err != nil {
			klog.Errorf("error stopping aite: %v", err)
		}
	}()
	apb.RegisterAiteServer(serv, as)

	reflection.Register(serv)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		klog.Exitf("cannot start listening, err: %v", err)
	}

	klog.Infof("aite server listening on %s", lis.Addr().String())
	defer serv.Stop()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := serv.Serve(lis); err != nil {
			klog.Exitf("cannot start server, %v", err)
		}
	}()
	<-done
}
