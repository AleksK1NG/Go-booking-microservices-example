package grpc_client

import (
	"context"

	"google.golang.org/grpc"

	"github.com/AleksK1NG/hotels-mocroservices/user/config"
	"github.com/AleksK1NG/hotels-mocroservices/user/internal/interceptors"
	"github.com/AleksK1NG/hotels-mocroservices/user/pkg/logger"
)

// NewSessionServiceClient
func NewSessionServiceConn(ctx context.Context, cfg *config.Config, logger logger.Logger) (*grpc.ClientConn, error) {
	sessGRPCConn, err := grpc.DialContext(ctx, cfg.GRPCServer.SessionGrpcServicePort,
		grpc.WithUnaryInterceptor(interceptors.GetInterceptor(logger)),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	return sessGRPCConn, nil
}
