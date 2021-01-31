package usecase

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/streadway/amqp"

	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/internal/image"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/pkg/logger"
)

type ImageUseCase struct {
	pgRepo  image.PgRepository
	awsRepo image.AWSRepository
	logger  logger.Logger
}

func NewImageUseCase(pgRepo image.PgRepository, awsRepo image.AWSRepository, logger logger.Logger) *ImageUseCase {
	return &ImageUseCase{pgRepo: pgRepo, awsRepo: awsRepo, logger: logger}
}

func (i *ImageUseCase) ResizeImage(ctx context.Context, delivery amqp.Delivery) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ImageUseCase.ResizeImage")
	defer span.Finish()

	i.logger.Infof("amqp.Delivery: %-v", delivery)

	return nil
}
