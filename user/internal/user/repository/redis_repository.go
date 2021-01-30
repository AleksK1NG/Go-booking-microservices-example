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

	"github.com/AleksK1NG/hotels-mocroservices/user/internal/models"
)

// UserRedisRepository
type UserRedisRepository struct {
	redisConn  *redis.Client
	prefix     string
	expiration time.Duration
}

// NewUserRedisRepository
func NewUserRedisRepository(redisConn *redis.Client, prefix string, expiration time.Duration) *UserRedisRepository {
	return &UserRedisRepository{redisConn: redisConn, prefix: prefix, expiration: expiration}
}

// SaveUser
func (u *UserRedisRepository) SaveUser(ctx context.Context, user *models.UserResponse) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRedisRepository.SaveUser")
	defer span.Finish()

	userBytes, err := json.Marshal(user)
	if err != nil {
		return errors.Wrap(err, "UserRedisRepository.SaveUser.json.Marshal")
	}

	if err := u.redisConn.SetEX(ctx, u.createKey(user.UserID), string(userBytes), u.expiration).Err(); err != nil {
		return errors.Wrap(err, "UserRedisRepository.SaveUser.redisConn.SetEX")
	}

	return nil
}

// GetUserByID
func (u *UserRedisRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.UserResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRedisRepository.GetUserByID")
	defer span.Finish()

	result, err := u.redisConn.Get(ctx, u.createKey(userID)).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "UserRedisRepository.GetUserByID.redisConn.Get")
	}

	var res models.UserResponse
	if err := json.Unmarshal(result, &res); err != nil {
		return nil, errors.Wrap(err, "UserRedisRepository.GetUserByID.json.Unmarshal")
	}
	return &res, nil
}

// DeleteUser
func (u *UserRedisRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRedisRepository.DeleteUser")
	defer span.Finish()

	if err := u.redisConn.Del(ctx, u.createKey(userID)).Err(); err != nil {
		return errors.Wrap(err, "UserRedisRepository.GetUserByID.redisConn.Del")
	}

	return nil
}

func (u *UserRedisRepository) createKey(userID uuid.UUID) string {
	return fmt.Sprintf("%s: %s", u.prefix, userID)
}
