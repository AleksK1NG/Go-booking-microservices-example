package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/AleksK1NG/hotels-mocroservices/comments/internal/models"
)

// CommPGRepo
type commPGRepo struct {
	db *pgxpool.Pool
}

// NewCommPGRepo
func NewCommPGRepo(db *pgxpool.Pool) *commPGRepo {
	return &commPGRepo{db: db}
}

// Create
func (c *commPGRepo) Create(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	panic("implement me")
}
