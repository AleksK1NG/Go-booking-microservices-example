package delivery

import (
	"context"

	"github.com/AleksK1NG/hotels-mocroservices/sessions/proto"
)

// SessionsService
type SessionsService struct {
}

func NewSessionGRPCDelivery() *SessionsService {
	return &SessionsService{}
}

// CreateSession
func (s *SessionsService) CreateSession(ctx context.Context, id *sessionService.UserID) (*sessionService.SessionID, error) {
	panic("implement me")
}

// GetIDBySession
func (s *SessionsService) GetIDBySession(ctx context.Context, id *sessionService.SessionID) (*sessionService.UserID, error) {
	panic("implement me")
}

// DeleteSession
func (s *SessionsService) DeleteSession(ctx context.Context, id *sessionService.SessionID) (*sessionService.Empty, error) {
	panic("implement me")
}

// CreateCsrfToken
func (s *SessionsService) CreateCsrfToken(ctx context.Context, input *sessionService.CsrfTokenInput) (*sessionService.CsrfToken, error) {
	panic("implement me")
}

// CheckCsrfToken
func (s *SessionsService) CheckCsrfToken(ctx context.Context, check *sessionService.CsrfTokenCheck) (*sessionService.CheckResult, error) {
	panic("implement me")
}
