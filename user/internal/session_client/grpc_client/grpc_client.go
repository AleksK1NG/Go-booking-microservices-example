package grpc_client

import (
	"context"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	traceutils "github.com/opentracing-contrib/go-grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/AleksK1NG/hotels-mocroservices/user/config"
	"github.com/AleksK1NG/hotels-mocroservices/user/internal/interceptors"
)

const (
	backoffLinear = 100 * time.Millisecond
)

// NewSessionServiceClient
func NewSessionServiceConn(ctx context.Context, cfg *config.Config, manager *interceptors.InterceptorManager) (*grpc.ClientConn, error) {
	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(backoffLinear)),
		grpc_retry.WithCodes(codes.NotFound, codes.Aborted),
	}

	sessGRPCConn, err := grpc.DialContext(
		ctx,
		cfg.GRPCServer.SessionGrpcServicePort,
		grpc.WithUnaryInterceptor(traceutils.OpenTracingClientInterceptor(manager.GetTracer())),
		grpc.WithUnaryInterceptor(manager.GetInterceptor()),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)),
	)
	if err != nil {
		return nil, err
	}

	return sessGRPCConn, nil
}
