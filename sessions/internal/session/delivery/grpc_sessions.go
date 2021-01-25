package delivery

import (
	"context"

	"github.com/opentracing/opentracing-go"

	"github.com/AleksK1NG/hotels-mocroservices/sessions/internal/session"
	"github.com/AleksK1NG/hotels-mocroservices/sessions/pkg/logger"
	sessionService "github.com/AleksK1NG/hotels-mocroservices/sessions/proto"
)

// SessionsService
type SessionsService struct {
	logger logger.Logger
	sessUC session.SessUseCase
}

func NewSessionsService(logger logger.Logger, sessUC session.SessUseCase) *SessionsService {
	return &SessionsService{logger: logger, sessUC: sessUC}
}

func (s *SessionsService) CreateSession(ctx context.Context, r *sessionService.CreateSessionRequest) (*sessionService.CreateSessionResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionsService.CreateSession")
	defer span.Finish()
	return nil, nil
}

func (s *SessionsService) GetSessionByID(ctx context.Context, r *sessionService.GetSessionByIDRequest) (*sessionService.GetSessionByIDResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionsService.GetSessionByID")
	defer span.Finish()
	return nil, nil
}

func (s *SessionsService) DeleteSession(ctx context.Context, r *sessionService.DeleteSessionRequest) (*sessionService.DeleteSessionResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionsService.DeleteSession")
	defer span.Finish()
	return nil, nil
}

func (s *SessionsService) CreateCsrfToken(ctx context.Context, r *sessionService.CreateCsrfTokenRequest) (*sessionService.CreateCsrfTokenResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionsService.CreateCsrfToken")
	defer span.Finish()
	return nil, nil
}

func (s *SessionsService) CheckCsrfToken(ctx context.Context, r *sessionService.CheckCsrfTokenRequest) (*sessionService.CheckCsrfTokenResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SessionsService.CheckCsrfToken")
	defer span.Finish()
	return nil, nil
}
