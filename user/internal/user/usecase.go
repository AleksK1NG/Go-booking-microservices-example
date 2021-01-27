package user

import (
	"context"

	"github.com/AleksK1NG/hotels-mocroservices/user/internal/models"
)

// UseCase
type UseCase interface {
	Register(ctx context.Context, user *models.User) (*models.UserResponse, error)
}
