package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "commPGRepo.Create")
	defer span.Finish()

	createCommentQuery := `INSERT INTO comments (hotel_id, user_id, message, photos, rating) 
	VALUES ($1, $2, $3, $4, $5) RETURNING comment_id, hotel_id, user_id, message, photos, rating, created_at, updated_at`

	var comm models.Comment
	if err := c.db.QueryRow(
		ctx,
		createCommentQuery,
		comment.HotelID,
		comment.UserID,
		comment.Message,
		comment.Photos,
		comment.Rating,
	).Scan(
		&comm.CommentID,
		&comm.HotelID,
		&comm.UserID,
		&comm.Message,
		&comm.Photos,
		&comm.Rating,
		&comm.CreatedAt,
		&comm.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}

	return &comm, nil
}
