package user

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/user/internal/models"
)

// UseCase
type UseCase interface {
	Register(ctx context.Context, user *models.User) (*models.UserResponse, error)
	Login(ctx context.Context, login models.Login) (*models.User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*models.UserResponse, error)
	CreateSession(ctx context.Context, userID uuid.UUID) (string, error)
}
