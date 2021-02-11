package hotels

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/internal/models"
)

// UseCase
type UseCase interface {
	GetHotels(ctx context.Context, page, size int64) (*models.HotelsListRes, error)
	GetHotelByID(ctx context.Context, hotelID uuid.UUID) (*models.Hotel, error)
	UpdateHotel(ctx context.Context, hotel *models.Hotel) (*models.Hotel, error)
	CreateHotel(ctx context.Context, hotel *models.Hotel) (*models.Hotel, error)
	UploadImage(ctx context.Context, data []byte, contentType, hotelID string) error
}
