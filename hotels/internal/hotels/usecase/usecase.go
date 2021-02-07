package usecase

import (
	"context"

	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/hotels/internal/hotels"
	"github.com/AleksK1NG/hotels-mocroservices/hotels/internal/models"
	"github.com/AleksK1NG/hotels-mocroservices/hotels/pkg/logger"
	"github.com/AleksK1NG/hotels-mocroservices/hotels/pkg/utils"
)

// HotelsUC Hotels usecase
type HotelsUC struct {
	hotelsRepo hotels.PGRepository
	logger     logger.Logger
}

// NewHotelsUC constructor
func NewHotelsUC(hotelsRepo hotels.PGRepository, logger logger.Logger) *HotelsUC {
	return &HotelsUC{hotelsRepo: hotelsRepo, logger: logger}
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

func (h *HotelsUC) GetHotels(ctx context.Context, query *utils.PaginationQuery) ([]*models.Hotel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HotelsUC.CreateHotel")
	defer span.Finish()

	return h.hotelsRepo.GetHotels(ctx, query)
}
