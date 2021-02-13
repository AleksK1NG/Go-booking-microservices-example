package grpc_client

import (
	"context"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	traceutils "github.com/opentracing-contrib/go-grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/internal/interceptors"
)

const (
	backoffLinear = 100 * time.Millisecond
)

func NewGRPCClientServiceConn(ctx context.Context, manager *interceptors.InterceptorManager, target string) (*grpc.ClientConn, error) {
	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(backoffLinear)),
		grpc_retry.WithCodes(codes.NotFound, codes.Aborted),
	}

	clientGRPCConn, err := grpc.DialContext(
		ctx,
		target,
		grpc.WithUnaryInterceptor(traceutils.OpenTracingClientInterceptor(manager.GetTracer())),
		grpc.WithUnaryInterceptor(manager.GetInterceptor()),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)),
	)
	if err != nil {
		return nil, err
	}

	return clientGRPCConn, nil
}
