package hotels

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/hotels/internal/models"
)

// Hotels use case
type UseCase interface {
	CreateHotel(ctx context.Context, hotel *models.Hotel) (*models.Hotel, error)
	UpdateHotel(ctx context.Context, hotel *models.Hotel) (*models.Hotel, error)
	GetHotelByID(ctx context.Context, hotelID uuid.UUID) (*models.Hotel, error)
}
