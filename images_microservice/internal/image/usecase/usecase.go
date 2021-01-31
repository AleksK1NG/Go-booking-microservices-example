package usecase

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/internal/image"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/internal/image/publisher"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/pkg/logger"
)

type ImageUseCase struct {
	pgRepo    image.PgRepository
	awsRepo   image.AWSRepository
	logger    logger.Logger
	publisher publisher.Publisher
}

func NewImageUseCase(pgRepo image.PgRepository, awsRepo image.AWSRepository, logger logger.Logger, publisher publisher.Publisher) *ImageUseCase {
	return &ImageUseCase{pgRepo: pgRepo, awsRepo: awsRepo, logger: logger, publisher: publisher}
}

func (i *ImageUseCase) ResizeImage(ctx context.Context, delivery amqp.Delivery) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ImageUseCase.ResizeImage")
	defer span.Finish()

	i.logger.Infof("amqp.Delivery: %-v", delivery)

	if err := i.publisher.Publish("images", "uploaded", "image/jpeg", []byte(time.Now().String())); err != nil {
		return errors.Wrap(err, "ImageUseCase.ResizeImage.Publish")
	}

	return nil
}
