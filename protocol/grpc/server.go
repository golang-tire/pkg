package grpc

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/golang-tire/pkg/protocol"

	"google.golang.org/grpc"

	"github.com/golang-tire/pkg/log"
	"github.com/golang-tire/pkg/protocol/grpc/middleware"
)

// RunServer runs gRPC service to publish all services
func RunServer(ctx context.Context, port int, logger log.Logger, controllers ...protocol.Controller) error {
	address := fmt.Sprintf(":%d", port)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	// gRPC server startup options
	opts := []grpc.ServerOption{}
	// add middleware
	opts = middleware.AddLogging(logger, opts)
	// register service
	server := grpc.NewServer(opts...)
	for _, controller := range controllers {
		controller.RegisterGRPCServer(server)
	}
	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			logger.Error("shutting down gRPC server...")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()
	// start gRPC server
	logger.Info("starting gRPC server...")
	return server.Serve(listen)
}
