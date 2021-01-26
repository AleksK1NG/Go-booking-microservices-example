package csrf

import "context"

// CSRF UseCase
type UseCase interface {
	GetCSRFToken(ctx context.Context, sesID string) (string, error)
	ValidateCSRFToken(ctx context.Context, sesID string, token string) (bool, error)
}
