package session

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/sessions/internal/models"
)

// Session RedisRepository
type RedisRepository interface {
	CreateSession(ctx context.Context, userID uuid.UUID) (*models.Session, error)
	GetSessionByID(ctx context.Context, sessID string) (*models.Session, error)
	DeleteSession(ctx context.Context, sessID string) error
}
