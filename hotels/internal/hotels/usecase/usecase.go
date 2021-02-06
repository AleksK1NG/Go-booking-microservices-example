package usecase

import (
	"context"

	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/hotels/internal/hotels"
	"github.com/AleksK1NG/hotels-mocroservices/hotels/internal/models"
	"github.com/AleksK1NG/hotels-mocroservices/hotels/pkg/logger"
)

type HotelsUC struct {
	hotelsRepo hotels.PGRepository
	logger     logger.Logger
}

func NewHotelsUC(hotelsRepo hotels.PGRepository, logger logger.Logger) *HotelsUC {
	return &HotelsUC{hotelsRepo: hotelsRepo, logger: logger}
}

func (h *HotelsUC) CreateHotel(ctx context.Context, hotel *models.Hotel) (*models.Hotel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HotelsUC.CreateHotel")
	defer span.Finish()

	return h.hotelsRepo.CreateHotel(ctx, hotel)
}

func (h *HotelsUC) UpdateHotel(ctx context.Context, hotel *models.Hotel) (*models.Hotel, error) {
	panic("implement me")
}

func (h *HotelsUC) GetHotelByID(ctx context.Context, hotelID uuid.UUID) (*models.Hotel, error) {
	panic("implement me")
}
