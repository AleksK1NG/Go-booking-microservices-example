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
	GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error)
	GetCSRFToken(ctx context.Context, sessionID string) (string, error)
	DeleteSession(ctx context.Context, sessionID string) error
	Update(ctx context.Context, user *models.UserUpdate) (*models.UserResponse, error)
}
