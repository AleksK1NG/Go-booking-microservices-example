package grpc

import (
	"context"

	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/config"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/pkg/logger"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/proto/image"
)

type ImageService struct {
	cfg    *config.Config
	logger logger.Logger
}

func NewImageService(cfg *config.Config, logger logger.Logger) *ImageService {
	return &ImageService{cfg: cfg, logger: logger}
}

func (i *ImageService) GetImageByID(ctx context.Context, req *imageService.GetByIDRequest) (*imageService.GetByIDResponse, error) {
	panic("implement me")
}
