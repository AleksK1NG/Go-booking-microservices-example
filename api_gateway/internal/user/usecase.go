package user

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/internal/models"
)

// UseCase
type UseCase interface {
	GetByID(ctx context.Context, userUUID uuid.UUID) (*models.UserResponse, error)
	GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error)
}
