package usecase

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/internal/hotels"
	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/internal/models"
	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/pkg/logger"
	hotelsService "github.com/AleksK1NG/hotels-mocroservices/api-gateway/proto/hotels"
)

// HotelsUseCase
type hotelsUseCase struct {
	logger        logger.Logger
	hotelsService hotelsService.HotelsServiceClient
	hotelsRepo    hotels.RedisRepository
}

// NewHotelsUseCase
func NewHotelsUseCase(logger logger.Logger, hotelsService hotelsService.HotelsServiceClient, hotelsRepo hotels.RedisRepository) *hotelsUseCase {
	return &hotelsUseCase{logger: logger, hotelsService: hotelsService, hotelsRepo: hotelsRepo}
}

// GetHotelByID
func (h *hotelsUseCase) GetHotelByID(ctx context.Context, hotelID uuid.UUID) (*models.Hotel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HotelsUseCase.GetHotelByID")
	defer span.Finish()

	hotelByID, err := h.hotelsService.GetHotelByID(ctx, &hotelsService.GetByIDReq{HotelID: hotelID.String()})
	if err != nil {
		return nil, errors.Wrap(err, "hotelsService.GetHotelByID")
	}

	fromProto, err := models.HotelFromProto(hotelByID.GetHotel())
	if err != nil {
		return nil, errors.Wrap(err, "HotelFromProto")
	}

	return fromProto, nil
}

// UpdateHotel
func (h *hotelsUseCase) UpdateHotel(ctx context.Context, hotel *models.Hotel) (*models.Hotel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HotelsUseCase.UpdateHotel")
	defer span.Finish()

	hotelRes, err := h.hotelsService.UpdateHotel(ctx, &hotelsService.UpdateHotelReq{
		HotelID:       hotel.HotelID.String(),
		Name:          hotel.Name,
		Email:         hotel.Email,
		Country:       hotel.Country,
		City:          hotel.City,
		Description:   hotel.Description,
		Location:      hotel.Location,
		Rating:        hotel.Rating,
		Image:         *hotel.Image,
		Photos:        hotel.Photos,
		CommentsCount: int64(hotel.CommentsCount),
		Latitude:      *hotel.Latitude,
		Longitude:     *hotel.Longitude,
	})
	if err != nil {
		return nil, errors.Wrap(err, "hotelsService.UpdateHotel")
	}

	fromProto, err := models.HotelFromProto(hotelRes.GetHotel())
	if err != nil {
		return nil, errors.Wrap(err, "HotelFromProto")
	}

	return fromProto, nil
}

// CreateHotel
func (h *hotelsUseCase) CreateHotel(ctx context.Context, hotel *models.Hotel) (*models.Hotel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HotelsUseCase.CreateHotel")
	defer span.Finish()

	hotelRes, err := h.hotelsService.CreateHotel(ctx, &hotelsService.CreateHotelReq{
		Name:          hotel.Name,
		Email:         hotel.Email,
		Country:       hotel.Country,
		City:          hotel.City,
		Description:   hotel.Description,
		Location:      hotel.Location,
		Rating:        hotel.Rating,
		Image:         *hotel.Image,
		Photos:        hotel.Photos,
		CommentsCount: int64(hotel.CommentsCount),
		Latitude:      *hotel.Latitude,
		Longitude:     *hotel.Longitude,
	})
	if err != nil {
		return nil, errors.Wrap(err, "hotelsService.CreateHotel")
	}

	fromProto, err := models.HotelFromProto(hotelRes.GetHotel())
	if err != nil {
		return nil, errors.Wrap(err, "HotelFromProto")
	}

	return fromProto, nil
}

// UploadImage
func (h *hotelsUseCase) UploadImage(ctx context.Context, data []byte, contentType, hotelID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HotelsUseCase.UploadImage")
	defer span.Finish()

	_, err := h.hotelsService.UploadImage(ctx, &hotelsService.UploadImageReq{
		HotelID:     hotelID,
		Data:        data,
		ContentType: contentType,
	})
	if err != nil {
		return errors.Wrap(err, "hotelsService.UploadImage")
	}

	return nil
}

// GetHotelsGetHotels
func (h *hotelsUseCase) GetHotels(ctx context.Context, page, size int64) (*models.HotelsListRes, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HotelsUseCase.GetHotels")
	defer span.Finish()

	hotelsRes, err := h.hotelsService.GetHotels(ctx, &hotelsService.GetHotelsReq{
		Page: page,
		Size: size,
	})
	if err != nil {
		return nil, errors.Wrap(err, "hotelsService.GetHotels")
	}

	hotelsList := make([]*models.Hotel, 0, len(hotelsRes.Hotels))
	for _, v := range hotelsRes.Hotels {
		hotel, err := models.HotelFromProto(v)
		if err != nil {
			return nil, errors.Wrap(err, "HotelFromProto")
		}
		hotelsList = append(hotelsList, hotel)
	}

	return &models.HotelsListRes{
		TotalCount: hotelsRes.GetTotalCount(),
		TotalPages: hotelsRes.GetTotalPages(),
		Page:       hotelsRes.GetPage(),
		Size:       hotelsRes.GetSize(),
		HasMore:    hotelsRes.GetHasMore(),
		Hotels:     hotelsList,
	}, nil
}
