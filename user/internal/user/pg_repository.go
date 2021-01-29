package user

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/user/internal/models"
)

// PGRepository
type PGRepository interface {
	Create(ctx context.Context, user *models.User) (*models.UserResponse, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*models.UserResponse, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.UserUpdate) (*models.UserResponse, error)
}
