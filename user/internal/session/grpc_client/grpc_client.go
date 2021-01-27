package grpc_client

import (
	"context"

	"google.golang.org/grpc"

	"github.com/AleksK1NG/hotels-mocroservices/user/config"
)

// NewSessionServiceClient
func NewSessionServiceConn(ctx context.Context, cfg *config.Config) (*grpc.ClientConn, error) {
	sessGRPCConn, err := grpc.DialContext(ctx, cfg.GRPCServer.SessionGrpcServicePort,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	return sessGRPCConn, nil
}
