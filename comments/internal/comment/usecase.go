package comment

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/comments/internal/models"
	"github.com/AleksK1NG/hotels-mocroservices/comments/pkg/utils"
)

// UseCase
type UseCase interface {
	Create(ctx context.Context, comment *models.Comment) (*models.Comment, error)
	GetByID(ctx context.Context, commentID uuid.UUID) (*models.Comment, error)
	Update(ctx context.Context, comment *models.Comment) (*models.Comment, error)
	GetByHotelID(ctx context.Context, hotelID uuid.UUID, query *utils.Pagination) (*models.CommentsList, error)
}
