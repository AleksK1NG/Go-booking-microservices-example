package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// UserPGRepository
type UserPGRepository struct {
	db *pgxpool.Pool
}

// NewUserPGRepository
func NewUserPGRepository(db *pgxpool.Pool) *UserPGRepository {
	return &UserPGRepository{db: db}
}

func (u *UserPGRepository) Register(ctx context.Context) {
	panic("implement me")
}
