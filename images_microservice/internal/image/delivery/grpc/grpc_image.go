package grpc

import (
	"context"

	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/status"

	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/config"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/internal/image"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/pkg/grpc_errors"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/pkg/logger"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/proto/image"
)

type ImageService struct {
	cfg     *config.Config
	logger  logger.Logger
	imageUC image.UseCase
}

func NewImageService(cfg *config.Config, logger logger.Logger, imageUC image.UseCase) *ImageService {
	return &ImageService{cfg: cfg, logger: logger, imageUC: imageUC}
}

func (i *ImageService) GetImageByID(ctx context.Context, req *imageService.GetByIDRequest) (*imageService.GetByIDResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ImageService.GetImageByID")
	defer span.Finish()

	imageUUID, err := uuid.FromString(req.GetImageID())
	if err != nil {
		i.logger.Errorf("uuid.FromString: %v", err)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "uuid.FromString: %v", err)
	}

	imageByID, err := i.imageUC.GetImageByID(ctx, imageUUID)
	if err != nil {
		i.logger.Errorf("uuid.FromString: %v", err)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "mageUC.GetImageByID: %v", err)
	}

	return &imageService.GetByIDResponse{Image: imageByID.ToProto()}, nil
}
