package repository

import (
	"context"

	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/internal/models"
)

// HotelRedisRepo
type hotelRedisRepo struct {
	redis *redis.Client
}

// NewHotelRedisRepo
func NewHotelRedisRepo(redis *redis.Client) *hotelRedisRepo {
	return &hotelRedisRepo{redis: redis}
}

// GetHotelByID
func (h *hotelRedisRepo) GetHotelByID(ctx context.Context, hotelID uuid.UUID) (*models.Hotel, error) {
	panic("implement me")
}

// SetHotel
func (h *hotelRedisRepo) SetHotel(ctx context.Context, hotel *models.Hotel) error {
	panic("implement me")
}

// DeleteHotel
func (h *hotelRedisRepo) DeleteHotel(ctx context.Context, hotelID uuid.UUID) error {
	panic("implement me")
}
