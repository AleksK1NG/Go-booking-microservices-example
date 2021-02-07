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

	point := utils.GeneratePointToGeoFromFloat64(*hotel.Latitude, *hotel.Longitude)
	createHotelQuery := `INSERT INTO hotels (name, location, description, image, photos, coordinates, email, country, city, rating) 
	VALUES ($1, $2, $3, $4, $5, ST_GeomFromEWKT($6), $7, $8, $9, $10) RETURNING hotel_id, created_at, updated_at`

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

	log.Printf("request :%-v", hotel)

	return hotel, nil
}

// UpdateHotel update single hotel
func (h *HotelsPGRepository) UpdateHotel(ctx context.Context, hotel *models.Hotel) (*models.Hotel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HotelsPGRepository.UpdateHotel")
	defer span.Finish()

	updateHotelQuery := `UPDATE hotels 
		SET email = COALESCE(NULLIF('', $1), email), name = $2, location = $3, description = $4, 
	 	country = $5, city = $6, coordinates = ST_GeomFromEWKT($7)
		WHERE hotel_id = $8
	    RETURNING hotel_id, email, name, location, description, comments_count, 
       	country, city, ((coordinates::POINT)[0])::decimal, ((coordinates::POINT)[1])::decimal, rating, photos, image, created_at, updated_at`

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
	).Scan(&hotel.HotelID,
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

	log.Printf("HOTEL BY ID: %v", res)

	return &res, nil
}

// GetHotelByID get hotel by id
func (h *HotelsPGRepository) GetHotelByID(ctx context.Context, hotelID uuid.UUID) (*models.Hotel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HotelsPGRepository.GetHotelByID")
	defer span.Finish()

	getHotelByIDQuery := `SELECT hotel_id, email, name, location, description, comments_count, 
       	country, city, ((coordinates::POINT)[0])::decimal, ((coordinates::POINT)[1])::decimal, rating, photos, image, created_at, updated_at 
		FROM hotels WHERE hotel_id = $1`

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

	log.Printf("HOTEL BY ID: %v", hotel)

	return &hotel, nil
}
