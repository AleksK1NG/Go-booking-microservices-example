package hotels

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/internal/models"
)

// RedisRepository
type RedisRepository interface {
	GetHotelByID(ctx context.Context, hotelID uuid.UUID) (*models.Hotel, error)
	SetHotel(ctx context.Context, hotel *models.Hotel) error
	DeleteHotel(ctx context.Context, hotelID uuid.UUID) error
}
