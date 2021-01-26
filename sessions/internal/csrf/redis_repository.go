package csrf

import "context"

// CSRF RedisRepository
type RedisRepository interface {
	Create(ctx context.Context, token string) error
	Check(ctx context.Context, token string) error
}
