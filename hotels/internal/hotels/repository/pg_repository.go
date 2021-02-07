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

// HotelsPGRepository
type HotelsPGRepository struct {
	db *pgxpool.Pool
}

// NewHotelsPGRepository
func NewHotelsPGRepository(db *pgxpool.Pool) *HotelsPGRepository {
	return &HotelsPGRepository{db: db}
}

// CreateHotel
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

// GetHotels
func (h *HotelsPGRepository) GetHotels(ctx context.Context, query *utils.PaginationQuery) ([]*models.Hotel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "HotelsPGRepository.GetHotels")
	defer span.Finish()

	getHotelsQuery := `SELECT hotel_id, email, name, location, description, comments_count, 
       	country, city, ((coordinates::POINT)[0])::decimal, ((coordinates::POINT)[1])::decimal, rating, photos, image, created_at, updated_at 
       	FROM hotels OFFSET $1 LIMIT $2`

	rows, err := h.db.Query(ctx, getHotelsQuery, query.GetOffset(), query.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "db.Query")
	}
	defer rows.Close()

	hotels := make([]*models.Hotel, 0, query.GetLimit())
	for rows.Next() {
		var hotel models.Hotel
		if err := rows.Scan(
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
			return nil, errors.Wrap(err, "rows.Scan")
		}
		hotels = append(hotels, &hotel)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows.Err")
	}

	log.Printf("HOTELS: %-v", hotels)

	return hotels, nil
}
