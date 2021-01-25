package delivery

import (
	"context"

	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/status"

	"github.com/AleksK1NG/hotels-mocroservices/sessions/internal/session"
	"github.com/AleksK1NG/hotels-mocroservices/sessions/pkg/grpc_errors"
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

	userUUID, err := uuid.FromString(r.UserID)
	if err != nil {
		s.logger.Errorf("uuid.FromString: %v", err)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "uuid.FromString: %v", err)
	}
	sess, err := s.sessUC.CreateSession(ctx, userUUID)
	if err != nil {
		s.logger.Errorf("sessUC.CreateSession: %v", err)
		return nil, status.Errorf(grpc_errors.ParseGRPCErrStatusCode(err), "sessUC.CreateSession: %v", err)
	}

	return &sessionService.CreateSessionResponse{Session: &sessionService.Session{
		UserID:    sess.UserID.String(),
		SessionID: sess.SessionID,
	}}, nil
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
