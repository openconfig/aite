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
	defer as.Stop()
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

	go serv.Serve(lis)
	<-done
}
