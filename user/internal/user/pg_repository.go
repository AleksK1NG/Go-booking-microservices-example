package user

import "context"

// PGRepository
type PGRepository interface {
	Register(ctx context.Context)
}
