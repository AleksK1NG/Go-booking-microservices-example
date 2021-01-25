package usecase

import (
	"context"

	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/sessions/internal/models"
	"github.com/AleksK1NG/hotels-mocroservices/sessions/internal/session"
)

type SessionUseCase struct {
	sessRepo session.RedisRepository
}

func NewSessionUseCase(sessRepo session.RedisRepository) *SessionUseCase {
	return &SessionUseCase{sessRepo: sessRepo}
}

func (s *SessionUseCase) CreateSession(ctx context.Context, userID uuid.UUID) (*models.Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionUseCase.CreateSession")
	defer span.Finish()
	return s.sessRepo.CreateSession(ctx, userID)
}

func (s *SessionUseCase) GetSessionByID(ctx context.Context, sessID string) (*models.Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionUseCase.GetSessionByID")
	defer span.Finish()
	return s.sessRepo.GetSessionByID(ctx, sessID)
}

func (s *SessionUseCase) DeleteSession(ctx context.Context, sessID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DeleteSession.GetSessionByID")
	defer span.Finish()
	return s.sessRepo.DeleteSession(ctx, sessID)
}
