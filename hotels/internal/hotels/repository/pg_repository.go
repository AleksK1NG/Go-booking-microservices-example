package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/hotels/internal/models"
	"github.com/AleksK1NG/hotels-mocroservices/hotels/pkg/utils"
)

type HotelsPGRepository struct {
	db *pgxpool.Pool
}

func NewHotelsPGRepository(db *pgxpool.Pool) *HotelsPGRepository {
	return &HotelsPGRepository{db: db}
}

func (h *HotelsPGRepository) CreateHotel(ctx context.Context, hotel *models.Hotel) (*models.Hotel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HotelsPGRepository.CreateHotel")
	defer span.Finish()

	point := utils.GeneratePointToGeoFromFloat64(hotel.Latitude.Float64, hotel.Longitude.Float64)
	createHotelQuery := `INSERT INTO hotels (name, location, description, image, photos, coordinates, email, country, city, rating) 
	VALUES ($1, $2, $3, $4, $5, ST_GeomFromEWKT($6), $7, $8, $9, $10) RETURNING hotel_id, created_at, updated_at`

	var res models.Hotel
	if err := h.db.QueryRow(
		ctx,
		createHotelQuery,
		hotel.Name,
		hotel.Location,
		hotel.Description,
		hotel.Image.String,
		hotel.Photos,
		point,
		hotel.Email,
		hotel.Country,
		hotel.City,
		hotel.Rating,
	).Scan(&res.HotelID, &res.CreatedAt, &res.UpdatedAt); err != nil {
		return nil, errors.Wrap(err, "CreateHotel.Scan")
	}

	hotel.HotelID = res.HotelID
	hotel.CreatedAt = res.CreatedAt
	hotel.UpdatedAt = res.UpdatedAt

	log.Printf("request :%-v", hotel)

	return hotel, nil
}

func (h *HotelsPGRepository) UpdateHotel(ctx context.Context, hotel *models.Hotel) (*models.Hotel, error) {
	panic("implement me")
}

func (h *HotelsPGRepository) GetHotelByID(ctx context.Context, hotelID uuid.UUID) (*models.Hotel, error) {
	panic("implement me")
}
