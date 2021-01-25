package repository

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/sessions/internal/models"
)

// SessionRedisRepo
type SessionRedisRepo struct {
	redis *redis.Client
}

func NewSessionRedisRepo(redis *redis.Client) *SessionRedisRepo {
	return &SessionRedisRepo{redis: redis}
}

func (s *SessionRedisRepo) CreateSession(ctx context.Context, userID uuid.UUID) (*models.Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionRedisRepo.CreateSession")
	defer span.Finish()
	panic("implement me")
}

func (s *SessionRedisRepo) GetSessionByID(ctx context.Context, sessID string) (*models.Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionRedisRepo.GetSessionByID")
	defer span.Finish()
	panic("implement me")
}

func (s *SessionRedisRepo) DeleteSession(ctx context.Context, sessID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionRedisRepo.DeleteSession")
	defer span.Finish()
	panic("implement me")
}
