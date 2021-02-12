package comments

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/internal/models"
)

// UseCase
type UseCase interface {
	CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error)
	GetCommByID(ctx context.Context, commentID uuid.UUID) (*models.Comment, error)
	UpdateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error)
	GetByHotelID(ctx context.Context, hotelID uuid.UUID, page, size int64) (*models.CommentsList, error)
}
