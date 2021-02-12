package comments

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/internal/models"
)

// RedisRepository
type RedisRepository interface {
	GetCommentByID(ctx context.Context, commentID uuid.UUID) (*models.Comment, error)
	SetComment(ctx context.Context, comment *models.Comment) error
	DeleteComment(ctx context.Context, commentID uuid.UUID) error
}
