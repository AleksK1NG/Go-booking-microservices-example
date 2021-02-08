package usecase

import (
	"context"
	"encoding/json"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"

	"github.com/AleksK1NG/hotels-mocroservices/hotels/internal/hotels"
	"github.com/AleksK1NG/hotels-mocroservices/hotels/internal/hotels/delivery/rabbitmq"
	"github.com/AleksK1NG/hotels-mocroservices/hotels/internal/models"
	"github.com/AleksK1NG/hotels-mocroservices/hotels/pkg/hotels_errors"
	"github.com/AleksK1NG/hotels-mocroservices/hotels/pkg/logger"
	"github.com/AleksK1NG/hotels-mocroservices/hotels/pkg/utils"
)

const (
	hotelIDHeader = "hotel_uuid"

	imagesExchange             = "images"
	uploadHotelImageRoutingKey = "upload_hotel_image"
)

// HotelsUC Hotels usecase
type HotelsUC struct {
	hotelsRepo    hotels.PGRepository
	logger        logger.Logger
	amqpPublisher rabbitmq.Publisher
}

// NewHotelsUC constructor
func NewHotelsUC(hotelsRepo hotels.PGRepository, logger logger.Logger, amqpPublisher rabbitmq.Publisher) *HotelsUC {
	return &HotelsUC{hotelsRepo: hotelsRepo, logger: logger, amqpPublisher: amqpPublisher}
}

// CreateHotel create new hotel
func (h *HotelsUC) CreateHotel(ctx context.Context, hotel *models.Hotel) (*models.Hotel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HotelsUC.CreateHotel")
	defer span.Finish()

	return h.hotelsRepo.CreateHotel(ctx, hotel)
}

// UpdateHotel update existing hotel
func (h *HotelsUC) UpdateHotel(ctx context.Context, hotel *models.Hotel) (*models.Hotel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HotelsUC.UpdateHotel")
	defer span.Finish()

	return h.hotelsRepo.UpdateHotel(ctx, hotel)
}

// GetHotelByID get hotel by uuid
func (h *HotelsUC) GetHotelByID(ctx context.Context, hotelID uuid.UUID) (*models.Hotel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HotelsUC.GetHotelByID")
	defer span.Finish()

	return h.hotelsRepo.GetHotelByID(ctx, hotelID)
}

// GetHotels
func (h *HotelsUC) GetHotels(ctx context.Context, query *utils.PaginationQuery) (*models.HotelsList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HotelsUC.CreateHotel")
	defer span.Finish()

	return h.hotelsRepo.GetHotels(ctx, query)
}

// UploadImage
func (h *HotelsUC) UploadImage(ctx context.Context, msg *models.UploadHotelImageMsg) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HotelsUC.UploadImage")
	defer span.Finish()

	headers := make(amqp.Table, 1)
	headers[hotelIDHeader] = msg.HotelID.String()
	if err := h.amqpPublisher.Publish(
		ctx,
		imagesExchange,
		uploadHotelImageRoutingKey,
		msg.ContentType,
		headers,
		msg.Data,
	); err != nil {
		return errors.Wrap(err, "UpdateUploadedAvatar.Publish")
	}

	return nil
}

// UpdateHotelImage
func (h *HotelsUC) UpdateHotelImage(ctx context.Context, delivery amqp.Delivery) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HotelsUC.UpdateHotelImage")
	defer span.Finish()

	var msg models.UpdateHotelImageMsg
	if err := json.Unmarshal(delivery.Body, &msg); err != nil {
		return errors.Wrap(err, "UpdateHotelImage.json.Unmarshal")
	}

	if err := h.hotelsRepo.UpdateHotelImage(ctx, msg.HotelID, msg.Image); err != nil {
		return err
	}

	return nil
}

func (h *HotelsUC) validateDeliveryHeaders(delivery amqp.Delivery) (*uuid.UUID, error) {
	h.logger.Infof("amqp.Delivery header: %-v", delivery.Headers)

	hotelUUID, ok := delivery.Headers[hotelIDHeader]
	if !ok {
		return nil, hotels_errors.ErrInvalidDeliveryHeaders
	}
	hotelID, ok := hotelUUID.(string)
	if !ok {
		return nil, hotels_errors.ErrInvalidUUID
	}

	parsedUUID, err := uuid.FromString(hotelID)
	if err != nil {
		return nil, errors.Wrap(err, "uuid.FromString")
	}

	return &parsedUUID, nil
}
