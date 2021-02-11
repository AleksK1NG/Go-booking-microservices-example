package usecase

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

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
			return nil, err
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
