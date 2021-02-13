package usecase

import (
	"context"

	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/sessions/internal/models"
	"github.com/AleksK1NG/hotels-mocroservices/sessions/internal/session"
)

type sessionUseCase struct {
	sessRepo session.RedisRepository
}

func NewSessionUseCase(sessRepo session.RedisRepository) *sessionUseCase {
	return &sessionUseCase{sessRepo: sessRepo}
}

func (s *sessionUseCase) CreateSession(ctx context.Context, userID uuid.UUID) (*models.Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionUseCase.CreateSession")
	defer span.Finish()
	return s.sessRepo.CreateSession(ctx, userID)
}

func (s *sessionUseCase) GetSessionByID(ctx context.Context, sessID string) (*models.Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sessionUseCase.GetSessionByID")
	defer span.Finish()
	return s.sessRepo.GetSessionByID(ctx, sessID)
}

func (s *sessionUseCase) DeleteSession(ctx context.Context, sessID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DeleteSession.GetSessionByID")
	defer span.Finish()
	return s.sessRepo.DeleteSession(ctx, sessID)
}
