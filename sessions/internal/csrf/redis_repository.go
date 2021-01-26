package csrf

import "context"

// CSRF RedisRepository
type RedisRepository interface {
	Create(ctx context.Context, token string) error
	GetToken(ctx context.Context, token string) (string, error)
}
