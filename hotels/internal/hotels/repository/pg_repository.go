package repository

import (
	"context"

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

	point := utils.GeneratePointToGeoFromFloat64(*hotel.Latitude, *hotel.Longitude)

	var res models.Hotel
	if err := h.db.QueryRow(
		ctx,
		createHotelQuery,
		hotel.Name,
		hotel.Location,
		hotel.Description,
		hotel.Image,
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

	return hotel, nil
}

// UpdateHotel update single hotel
func (h *HotelsPGRepository) UpdateHotel(ctx context.Context, hotel *models.Hotel) (*models.Hotel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HotelsPGRepository.UpdateHotel")
	defer span.Finish()

	point := utils.GeneratePointToGeoFromFloat64(*hotel.Latitude, *hotel.Longitude)
	var res models.Hotel
	if err := h.db.QueryRow(
		ctx,
		updateHotelQuery,
		hotel.Email,
		hotel.Name,
		hotel.Location,
		hotel.Description,
		hotel.Country,
		hotel.City,
		point,
		hotel.HotelID,
	).Scan(
		&res.HotelID,
		&res.Email,
		&res.Name,
		&res.Location,
		&res.Description,
		&res.CommentsCount,
		&res.Country,
		&res.City,
		&res.Latitude,
		&res.Longitude,
		&res.Rating,
		&res.Photos,
		&res.Image,
		&res.CreatedAt,
		&res.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "db.QueryRow.Scan")
	}

	return &res, nil
}

// GetHotelByID get hotel by id
func (h *HotelsPGRepository) GetHotelByID(ctx context.Context, hotelID uuid.UUID) (*models.Hotel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HotelsPGRepository.GetHotelByID")
	defer span.Finish()

	var hotel models.Hotel
	if err := h.db.QueryRow(ctx, getHotelByIDQuery, hotelID).Scan(
		&hotel.HotelID,
		&hotel.Email,
		&hotel.Name,
		&hotel.Location,
		&hotel.Description,
		&hotel.CommentsCount,
		&hotel.Country,
		&hotel.City,
		&hotel.Latitude,
		&hotel.Longitude,
		&hotel.Rating,
		&hotel.Photos,
		&hotel.Image,
		&hotel.CreatedAt,
		&hotel.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "db.QueryRow.Scan")
	}

	return &hotel, nil
}
