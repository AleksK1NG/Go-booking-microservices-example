package user

import "context"

// UseCase
type UseCase interface {
	Register(ctx context.Context)
}
