package rest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/golang-tire/pkg/protocol"
	"github.com/golang-tire/pkg/protocol/rest/swagger"
	"github.com/rs/cors"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/golang-tire/pkg/log"
	"github.com/golang-tire/pkg/protocol/rest/middleware"
)

// RunServer runs HTTP/REST gateway
func RunServer(ctx context.Context, httpPort, grpcPort int, logger log.Logger, controllers ...protocol.Controller) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	normalMux := http.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	swaggerService := swagger.New(logger)

	grpcAddress := fmt.Sprintf(":%d", grpcPort)
	for _, controller := range controllers {
		err := controller.RegisterRestEndpoint(ctx, mux, grpcAddress, opts)
		if err != nil {
			logger.Errorf("failed to register service rest endpoint: %s", zap.String("reason", err.Error()))
		}
		swaggerService.Register(controller.ReturnSwaggerDefinitions())
	}

	// TODO swagger url and base url, shall to be customizable
	normalMux.HandleFunc("/v1/swagger/", swaggerService.Handler)
	normalMux.Handle("/", cors.AllowAll().Handler(mux))

	httpAddress := fmt.Sprintf(":%d", httpPort)
	srv := &http.Server{
		Addr: httpAddress,
		// add handler with middleware
		Handler: middleware.AddRequestID(middleware.AddLogger(logger, normalMux)),
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
		}
		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}()

	logger.Info("starting HTTP/REST gateway...")
	return srv.ListenAndServe()
}
