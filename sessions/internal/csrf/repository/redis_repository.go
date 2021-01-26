package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

// CsrfRepository
type CsrfRepository struct {
	redis    *redis.Client
	prefix   string
	duration int
}

// NewRedisRepository
func NewCsrfRepository(redis *redis.Client, prefix string, duration int) *CsrfRepository {
	return &CsrfRepository{redis: redis, prefix: prefix, duration: duration}
}

// Create csrf token
func (r *CsrfRepository) Create(ctx context.Context, token string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CsrfRepository.Create")
	defer span.Finish()

	if err := r.redis.SetEX(ctx, r.createKey(token), token, time.Duration(r.duration)*time.Minute).Err(); err != nil {
		return errors.Wrap(err, "CsrfRepository.Create.redis.SetEX")
	}

	return nil
}

// Check csrf token
func (r *CsrfRepository) GetToken(ctx context.Context, token string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CsrfRepository.Check")
	defer span.Finish()

	token, err := r.redis.Get(ctx, r.createKey(token)).Result()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (r *CsrfRepository) createKey(token string) string {
	return fmt.Sprintf("%s: %s", r.prefix, token)
}
