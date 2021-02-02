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
	"io/ioutil"
	"sync"

	"github.com/disintegration/gift"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"

	img "github.com/AleksK1NG/hotels-mocroservices/images-microservice/internal/images"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/internal/images/publisher"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/internal/models"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/pkg/image_errors"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/pkg/images"
	"github.com/AleksK1NG/hotels-mocroservices/images-microservice/pkg/logger"
)

const (
	userUUIDHeader        = "user_uuid"
	imagesExchange        = "images"
	resizeImageRoutingKey = "uploaded"
	resizeWidth           = 1024
	resizeHeight          = 0
	imageCreateRoutingKey = "image_create"
)

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
			gift.Resize(resizeWidth, resizeHeight, gift.LanczosResampling),
			gift.Contrast(20),
			gift.Brightness(7),
			gift.Gamma(0.5),
		)
	}}
	return &ImageUseCase{pgRepo: pgRepo, awsRepo: awsRepo, logger: logger, publisher: publisher, resizerPool: resizerPool}
}

func (i *ImageUseCase) Create(ctx context.Context, delivery amqp.Delivery) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ImageUseCase.Create")
	defer span.Finish()

	i.logger.Infof("amqp.Delivery: %-v", delivery.DeliveryTag)

	var msg models.UploadImageMsg
	if err := json.Unmarshal(delivery.Body, &msg); err != nil {
		return err
	}

	createdImage, err := i.pgRepo.Create(ctx, &models.Image{
		ImageID:    msg.ImageID,
		ImageURL:   msg.ImageURL,
		IsUploaded: msg.IsUploaded,
	})
	if err != nil {
		return err
	}

	msgBytes, err := json.Marshal(createdImage)
	if err != nil {
		return errors.Wrap(err, "ImageUseCase.Create.json.Marshal")
	}

	headers := make(amqp.Table)
	headers["user_uuid"] = delivery.Headers["user_uuid"]

	i.logger.Infof("Create PUBLISH USER ***************** %-v", headers)

	if err := i.publisher.Publish(
		ctx,
		"users",
		"update_avatar_key",
		delivery.ContentType,
		headers,
		msgBytes,
	); err != nil {
		return errors.Wrap(err, "ImageUseCase.Create.Publish")
	}

	return nil
}

func (i *ImageUseCase) ResizeImage(ctx context.Context, delivery amqp.Delivery) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ImageUseCase.ResizeImage")
	defer span.Finish()

	i.logger.Infof("amqp.Delivery: %-v", delivery.DeliveryTag)

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

	msg := &models.UploadImageMsg{
		UserID:     *parsedUUID,
		ImageURL:   "url",
		IsUploaded: true,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "ImageUseCase.ResizeImage.json.Marshal")
	}

	headers := make(amqp.Table)
	headers["user_uuid"] = delivery.Headers["user_uuid"]

	i.logger.Infof("PUBLISH USER ***************** %-v", headers)

	if err := i.publisher.Publish(
		ctx,
		"images",
		"create_image_key",
		delivery.ContentType,
		headers,
		msgBytes,
	); err != nil {
		return errors.Wrap(err, "ImageUseCase.ResizeImage.Publish")
	}

	// if err := i.publisher.Publish(
	// 	ctx,
	// 	resizeImageExchange,
	// 	resizeImageRoutingKey,
	// 	delivery.ContentType,
	// 	msgBytes,
	// ); err != nil {
	// 	return errors.Wrap(err, "ImageUseCase.ResizeImage.Publish")
	// }

	return nil
}

func (i *ImageUseCase) validateDeliveryHeaders(delivery amqp.Delivery) (*uuid.UUID, error) {
	i.logger.Infof("amqp.Delivery header: %-v", delivery.Headers)

	userUUID, ok := delivery.Headers[userUUIDHeader]
	if !ok {
		return nil, errors.Wrap(image_errors.ErrInvalidDeliveryHeaders, "ImageUseCase.ResizeImage.Publish")
	}
	userID, ok := userUUID.(string)
	if !ok {
		return nil, errors.Wrap(image_errors.ErrInvalidUUID, "ImageUseCase.ResizeImage.Publish")
	}

	parsedUUID, err := uuid.FromString(userID)
	if err != nil {
		return nil, errors.Wrap(err, "ImageUseCase.ResizeImage.uuid.FromString")
	}

	return &parsedUUID, nil
}

func (i *ImageUseCase) processImage(img []byte) ([]byte, error) {
	src, imageType, err := image.Decode(bytes.NewReader(img))
	if err != nil {
		return nil, err
	}

	imgResizer, ok := i.resizerPool.Get().(*images.ImgResizer)
	if !ok {
		return nil, errors.Wrap(image_errors.ErrInternalServerError, "ImageUseCase.ResizeImage.resizerPool.Get")
	}
	defer i.resizerPool.Put(imgResizer)
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
		return nil, image_errors.ErrInvalidImageFormat
	}

	return imgResizer.Buffer.Bytes(), nil
}

func (i *ImageUseCase) uploadToAWS(data []byte) error {
	// file, err := os.Create("image.jpeg")
	// if err != nil {
	// 	return err
	// }

	r := bufio.NewReader(bytes.NewReader(data))

	written, err := io.Copy(ioutil.Discard, r)
	if err != nil {
		return err
	}
	i.logger.Infof("written: %v", written)
	return nil
}
