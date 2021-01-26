package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

// RedisRepository
type RedisRepository struct {
	redis    *redis.Client
	prefix   string
	duration int
}

// NewRedisRepository
func NewRedisRepository(redis *redis.Client, prefix string, duration int) *RedisRepository {
	return &RedisRepository{redis: redis, prefix: prefix, duration: duration}
}

// Create csrf token
func (r *RedisRepository) Create(ctx context.Context, token string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisRepository.Create")
	defer span.Finish()

	if err := r.redis.SetEX(ctx, r.createKey(token), token, time.Duration(r.duration)*time.Minute).Err(); err != nil {
		return errors.Wrap(err, "RedisRepository.Create.redis.SetEX")
	}

	return nil
}

// Check csrf token
func (r *RedisRepository) Check(ctx context.Context, token string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RedisRepository.Check")
	defer span.Finish()

	_, err := r.redis.Get(ctx, r.createKey(token)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}

	return errors.New("token is not valid")
}

func (r *RedisRepository) createKey(token string) string {
	return fmt.Sprintf("%s: %s", r.prefix, token)
}
