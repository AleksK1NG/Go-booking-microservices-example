package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/internal/models"
)

const (
	prefix     = "comments:"
	expiration = time.Second * 3600
)

// HotelRedisRepo
type hotelRedisRepo struct {
	redisConn *redis.Client
}

// NewHotelRedisRepo
func NewHotelRedisRepo(redisConn *redis.Client) *hotelRedisRepo {
	return &hotelRedisRepo{redisConn: redisConn}
}

// GetHotelByID
func (h *hotelRedisRepo) GetHotelByID(ctx context.Context, hotelID uuid.UUID) (*models.Hotel, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "hotelRedisRepo.GetHotelByID")
	defer span.Finish()

	result, err := h.redisConn.Get(ctx, h.createKey(hotelID)).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "hotelRedisRepo.GetHotelByID")
	}

	var res models.Hotel
	if err := json.Unmarshal(result, &res); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}

	return &res, nil
}

// SetHotel
func (h *hotelRedisRepo) SetHotel(ctx context.Context, hotel *models.Hotel) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "hotelRedisRepo.SetHotel")
	defer span.Finish()

	hotelBytes, err := json.Marshal(hotel)
	if err != nil {
		return errors.Wrap(err, "hotelRedisRepo.Marshal")
	}

	if err := h.redisConn.SetEX(ctx, h.createKey(hotel.HotelID), string(hotelBytes), expiration).Err(); err != nil {
		return errors.Wrap(err, "hotelRedisRepo.SetEX")
	}

	return nil
}

// DeleteHotel
func (h *hotelRedisRepo) DeleteHotel(ctx context.Context, hotelID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "hotelRedisRepo.DeleteHotel")
	defer span.Finish()

	if err := h.redisConn.Del(ctx, h.createKey(hotelID)).Err(); err != nil {
		return errors.Wrap(err, "hotelRedisRepo.DeleteHotel.Del")
	}

	return nil
}

func (h *hotelRedisRepo) createKey(hotelID uuid.UUID) string {
	return fmt.Sprintf("%s: %s", prefix, hotelID.String())
}
