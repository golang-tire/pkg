package grpcgw

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"google.golang.org/grpc/reflection"

	"github.com/golang-tire/pkg/log"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcCtxTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
)

type serverOptions struct {
	httpPort        int
	grpcPort        int
	swaggerBaseURL  string
	serveMuxOptions []runtime.ServeMuxOption
}

// A ServerOption sets options such as ports, paths parameters, etc.
type ServerOption interface {
	apply(*serverOptions)
}

// funcServerOption wraps a function that modifies serverOptions into an
// implementation of the Options interface.
type funcServerOption struct {
	f func(*serverOptions)
}

func (fdo *funcServerOption) apply(do *serverOptions) {
	fdo.f(do)
}

// Controller interface provide a way for each grpc and rest service to register their services
// in server
type Controller interface {
	InitRest(ctx context.Context, conn *grpc.ClientConn, mux *runtime.ServeMux)
	InitGrpc(ctx context.Context, server *grpc.Server)
}

// Interceptor interface provide a way to add custom interceptors for Unary and Stream server
type Interceptor struct {
	Unary  grpc.UnaryServerInterceptor
	Stream grpc.StreamServerInterceptor
}

var (
	controllers  []Controller
	interceptors []Interceptor
	lock         sync.RWMutex

	httpAddr string
	grpcAddr string

	defaultServerOptions = serverOptions{
		httpPort:       8080,
		grpcPort:       9090,
		swaggerBaseURL: "/v1/swagger",
	}
)

func newFuncServerOption(f func(*serverOptions)) *funcServerOption {
	return &funcServerOption{
		f: f,
	}
}

// HttpPort returns a ServerOption that will apply httpPort option
func HttpPort(n int) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.httpPort = n
	})
}

// GrpcPort returns a ServerOption that will apply grpcPort option
func GrpcPort(n int) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.grpcPort = n
	})
}

// SwaggerBaseURL returns a ServerOption that will apply swaggerBaseURL option
func SwaggerBaseURL(s string) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.swaggerBaseURL = s
	})
}

// ServeMuxOptions returns a ServerOption that will apply ServeMuxOptions option for mux
func ServeMuxOptions(opts ...runtime.ServeMuxOption) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.serveMuxOptions = opts
	})
}

// RegisterController register a controller
func RegisterController(c Controller) {
	lock.Lock()
	defer lock.Unlock()
	controllers = append(controllers, c)
}

// RegisterInterceptors register custom interceptors
func RegisterInterceptors(i Interceptor) {
	lock.Lock()
	defer lock.Unlock()
	interceptors = append(interceptors, i)
}

func newGrpcServer() *grpc.Server {
	unaryMiddlewares := []grpc.UnaryServerInterceptor{
		grpcRecovery.UnaryServerInterceptor(),
		grpcCtxTags.UnaryServerInterceptor(),
		grpcZap.UnaryServerInterceptor(log.Logger()),
	}

	streamMiddlewares := []grpc.StreamServerInterceptor{
		grpcRecovery.StreamServerInterceptor(),
		grpcCtxTags.StreamServerInterceptor(),
		grpcZap.StreamServerInterceptor(log.Logger()),
	}

	for i := range interceptors {
		if interceptors[i].Unary != nil {
			unaryMiddlewares = append(unaryMiddlewares, interceptors[i].Unary)
		}
		if interceptors[i].Stream != nil {
			streamMiddlewares = append(streamMiddlewares, interceptors[i].Stream)
		}
	}
	c := grpc.NewServer(
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(unaryMiddlewares...)),
		grpc.StreamInterceptor(grpcMiddleware.ChainStreamServer(streamMiddlewares...)),
	)
	reflection.Register(c)

	return c
}

// gRPCClient creates a new GRPC client conn
func gRPCClient() (*grpc.ClientConn, error) {
	return grpc.Dial(grpcAddr, grpc.WithInsecure())
}

// Serve start the server and wait
func serveHTTP(ctx context.Context, opts serverOptions) (func() error, error) {

	var (
		normalMux = http.NewServeMux()
		mux       = runtime.NewServeMux(opts.serveMuxOptions...)
	)
	c, err := gRPCClient()
	if err != nil {
		return nil, err
	}

	sw := &swaggerServer{swaggerBaseURL: opts.swaggerBaseURL}
	normalMux.HandleFunc(opts.swaggerBaseURL, sw.swaggerHandler)
	for i := range controllers {
		controllers[i].InitRest(ctx, c, mux)
	}

	normalMux.Handle("/", cors.AllowAll().Handler(mux))
	srv := http.Server{
		Addr:    httpAddr,
		Handler: normalMux,
	}
	log.Info("start http on", log.Any("address", httpAddr))
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return func() error {
		nCtx, cnl := context.WithTimeout(context.Background(), time.Second)
		defer cnl()

		return srv.Shutdown(nCtx)
	}, nil
}

// Serve start the server and wait
func serveGRPC(ctx context.Context) (func() error, error) {
	srv := newGrpcServer()
	for i := range controllers {
		controllers[i].InitGrpc(ctx, srv)
	}

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return nil, err
	}
	log.Info("start grpc on", log.Any("address", grpcAddr))
	go func() {
		err := srv.Serve(lis)
		if err != nil {
			log.Error("Connection Closed", log.Err(err))
		}
	}()

	return lis.Close, nil
}

// Serve will start http and grpc server on given ports
// and will use ctx.Done() channel to stop servers
func Serve(ctx context.Context, opt ...ServerOption) error {

	opts := defaultServerOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	httpAddr = fmt.Sprintf(":%d", opts.httpPort)
	grpcAddr = fmt.Sprintf(":%d", opts.grpcPort)

	lock.RLock()
	defer lock.RUnlock()

	grpcFn, err := serveGRPC(ctx)
	if err != nil {
		return err
	}
	httpFn, err := serveHTTP(ctx, opts)
	if err != nil {
		return err
	}

	<-ctx.Done()
	e1 := httpFn()
	e2 := grpcFn()

	if e1 != nil {
		return e1
	}

	return e2
}
