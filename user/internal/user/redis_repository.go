package user

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/user/internal/models"
)

// RedisRepository
type RedisRepository interface {
	SaveUser(ctx context.Context, user *models.UserResponse) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.UserResponse, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}
