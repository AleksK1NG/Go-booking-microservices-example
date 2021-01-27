package user

import (
	"context"

	"github.com/AleksK1NG/hotels-mocroservices/user/internal/models"
)

// PGRepository
type PGRepository interface {
	Create(ctx context.Context, user *models.User) (*models.UserResponse, error)
}
