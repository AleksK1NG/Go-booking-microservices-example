package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/hotels/internal/models"
)

type HotelsPGRepository struct {
	db *pgxpool.Pool
}

func NewHotelsPGRepository(db *pgxpool.Pool) *HotelsPGRepository {
	return &HotelsPGRepository{db: db}
}

func (h *HotelsPGRepository) CreateHotel(ctx context.Context, hotel *models.Hotel) (*models.Hotel, error) {
	panic("implement me")
}

func (h *HotelsPGRepository) UpdateHotel(ctx context.Context, hotel *models.Hotel) (*models.Hotel, error) {
	panic("implement me")
}

func (h *HotelsPGRepository) GetHotelByID(ctx context.Context, hotelID uuid.UUID) (*models.Hotel, error) {
	panic("implement me")
}
