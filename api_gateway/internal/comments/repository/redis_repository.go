package repository

import (
	"context"

	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/internal/models"
)

// CommRedisRepository
type commRedisRepository struct {
	redis *redis.Client
}

// NewCommRedisRepository
func NewCommRedisRepository(redis *redis.Client) *commRedisRepository {
	return &commRedisRepository{redis: redis}
}

func (c *commRedisRepository) GetCommentByID(ctx context.Context, commentID uuid.UUID) (*models.Comment, error) {
	panic("implement me")
}

func (c *commRedisRepository) SetComment(ctx context.Context, comment *models.Comment) error {
	panic("implement me")
}

func (c *commRedisRepository) DeleteComment(ctx context.Context, commentID uuid.UUID) error {
	panic("implement me")
}
