package usecase

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"sync"
	"time"

	"github.com/disintegration/gift"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"

	img "github.com/AleksK1NG/hotels-mocroservices/images-microservice/internal/images"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/internal/images/publisher"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/internal/models"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/pkg/images"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/pkg/logger"
)

// var resizerPool = &sync.Pool{New: func() interface{} {
// 	return images.NewImgResizer(
// 		gift.Resize(1024, 0, gift.LanczosResampling),
// 		gift.Contrast(20),
// 		gift.Brightness(7),
// 		gift.Gamma(0.5),
// 	)
// }}

type ImageUseCase struct {
	pgRepo      img.PgRepository
	awsRepo     img.AWSRepository
	logger      logger.Logger
	publisher   publisher.Publisher
	resizerPool *sync.Pool
}

func NewImageUseCase(pgRepo img.PgRepository, awsRepo img.AWSRepository, logger logger.Logger, publisher publisher.Publisher) *ImageUseCase {
	resizerPool := &sync.Pool{New: func() interface{} {
		return images.NewImgResizer(
			gift.Resize(1024, 0, gift.LanczosResampling),
			gift.Contrast(20),
			gift.Brightness(7),
			gift.Gamma(0.5),
		)
	}}
	return &ImageUseCase{pgRepo: pgRepo, awsRepo: awsRepo, logger: logger, publisher: publisher, resizerPool: resizerPool}
}

func (i *ImageUseCase) ResizeImage(ctx context.Context, delivery amqp.Delivery) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ImageUseCase.ResizeImage")
	defer span.Finish()

	i.logger.Infof("amqp.Delivery: %-v", delivery)

	parsedUUID, err := i.validateDeliveryHeaders(delivery)
	if err != nil {
		return err
	}

	processedImage, err := i.processImage(delivery.Body)
	if err != nil {
		return err
	}

	if err := i.uploadToAWS(processedImage); err != nil {
		return err
	}

	msg := &models.UploadedImageMsg{
		ImageID:    uuid.NewV4(),
		UserID:     *parsedUUID,
		ImageURL:   "url",
		IsUploaded: true,
		CreatedAt:  time.Now().UTC(),
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "ImageUseCase.ResizeImage.json.Marshal")
	}

	if err := i.publisher.Publish(
		ctx,
		"images",
		"uploaded",
		"image/jpeg",
		msgBytes,
	); err != nil {
		return errors.Wrap(err, "ImageUseCase.ResizeImage.Publish")
	}

	return nil
}

func (i *ImageUseCase) validateDeliveryHeaders(delivery amqp.Delivery) (*uuid.UUID, error) {
	i.logger.Infof("amqp.Delivery header: %-v", delivery.Headers)

	userUUID, ok := delivery.Headers["user_uuid"]
	if !ok {
		return nil, errors.Wrap(errors.New("Delivery header user_id is required"), "ImageUseCase.ResizeImage.Publish")
	}
	userID, ok := userUUID.(string)
	if !ok {
		return nil, errors.Wrap(errors.New("invalid user id"), "ImageUseCase.ResizeImage.Publish")
	}

	parsedUUID, err := uuid.FromString(userID)
	if err != nil {
		return nil, errors.Wrap(err, "ImageUseCase.ResizeImage.uuid.FromString")
	}
	i.logger.Infof("parsedUUID: %-v", parsedUUID.String())

	return &parsedUUID, nil
}

func (i *ImageUseCase) processImage(img []byte) ([]byte, error) {
	src, imageType, err := image.Decode(bytes.NewReader(img))
	if err != nil {
		return nil, err
	}

	imgResizer, ok := i.resizerPool.Get().(*images.ImgResizer)
	if !ok {
		return nil, errors.New("resizerPool.Get casting")
	}
	imgResizer.Buffer.Reset()

	dst := image.NewNRGBA(imgResizer.Gift.Bounds(src.Bounds()))
	imgResizer.Gift.Draw(dst, src)

	switch imageType {
	case "png":
		err = png.Encode(imgResizer.Buffer, dst)
		if err != nil {
			return nil, err
		}
	case "jpeg":
		err = jpeg.Encode(imgResizer.Buffer, dst, nil)
		if err != nil {
			return nil, err
		}
	case "jpg":
		err = jpeg.Encode(imgResizer.Buffer, dst, nil)
		if err != nil {
			return nil, err
		}
	case "gif":
		err = gif.Encode(imgResizer.Buffer, dst, nil)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("invalid image format")
	}

	return imgResizer.Buffer.Bytes(), nil
}

func (i *ImageUseCase) uploadToAWS(data []byte) error {
	file, err := os.Create("image.jpeg")
	if err != nil {
		return err
	}

	r := bufio.NewReader(bytes.NewReader(data))

	written, err := io.Copy(file, r)
	if err != nil {
		return err
	}
	i.logger.Infof("written: %v", written)
	return nil
}
