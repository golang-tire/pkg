package protocol

import (
	"context"

	"github.com/golang-tire/pkg/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

// Controller protocol controller object
type Controller interface {
	RegisterRestEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
	RegisterGRPCServer(grpcServer *grpc.Server)
	ReturnSwaggerDefinitions() (types.MapSI, types.MapSI, types.MapSI, bool)
}
