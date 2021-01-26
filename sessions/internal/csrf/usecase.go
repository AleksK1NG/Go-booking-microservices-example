package csrf

import "context"

// CSRF UseCase
type UseCase interface {
	CreateToken(ctx context.Context, sesID string, timeStamp int64) (string, error)
	CheckToken(ctx context.Context, sesID string, token string) (bool, error)
}
