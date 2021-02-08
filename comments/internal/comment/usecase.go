package comment

import (
	"context"

	"github.com/AleksK1NG/hotels-mocroservices/comments/internal/models"
)

// UseCase
type UseCase interface {
	Create(ctx context.Context, comment *models.Comment) (*models.Comment, error)
}
