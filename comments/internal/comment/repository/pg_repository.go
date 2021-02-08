package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/comments/internal/models"
	"github.com/AleksK1NG/hotels-mocroservices/comments/pkg/utils"
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

// GetByID
func (c *commPGRepo) GetByID(ctx context.Context, commentID uuid.UUID) (*models.Comment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commPGRepo.GetByID")
	defer span.Finish()

	getCommByIDQuery := `SELECT comment_id, hotel_id, user_id, message, photos, rating, created_at, updated_at FROM comments WHERE comment_id = $1`

	var comm models.Comment
	if err := c.db.QueryRow(ctx, getCommByIDQuery, commentID).Scan(
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

// Update
func (c *commPGRepo) Update(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commPGRepo.Update")
	defer span.Finish()

	updateCommentQuery := `UPDATE comments SET message = COALESCE(NULLIF($1, ''), message), rating = $2, photos = $3
	WHERE comment_id = $4
	RETURNING comment_id, hotel_id, user_id, message, photos, rating, created_at, updated_at`

	var comm models.Comment
	if err := c.db.QueryRow(ctx, updateCommentQuery, comment.Message, comment.Rating, comment.Photos, comment.CommentID).Scan(
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

// GetByHotelID
func (c *commPGRepo) GetByHotelID(ctx context.Context, hotelID uuid.UUID, query *utils.Pagination) (*models.CommentsList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "commPGRepo.GetByHotelID")
	defer span.Finish()

	getTotalCountQuery := `SELECT count(comment_id) as total FROM comments WHERE hotel_id = $1`
	var totalCount int
	if err := c.db.QueryRow(ctx, getTotalCountQuery, hotelID).Scan(&totalCount); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}

	if totalCount == 0 {
		return &models.CommentsList{
			TotalCount: 0,
			TotalPages: 0,
			Page:       0,
			Size:       0,
			HasMore:    false,
			Comments:   make([]*models.Comment, 0),
		}, nil
	}

	getCommentByHotelIDQuery := `SELECT comment_id, hotel_id, user_id, message, photos, rating, created_at, updated_at FROM comments
	WHERE hotel_id = $1 OFFSET $2 LIMIT $3`

	var commentsList []*models.Comment
	rows, err := c.db.Query(ctx, getCommentByHotelIDQuery, hotelID, query.GetOffset(), query.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "db.Query")
	}
	defer rows.Close()

	for rows.Next() {
		var comm models.Comment
		if err := rows.Scan(
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
		commentsList = append(commentsList, &comm)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows.Err")
	}

	return &models.CommentsList{
		TotalCount: totalCount,
		TotalPages: query.GetTotalPages(totalCount),
		Page:       query.GetPage(),
		Size:       query.GetSize(),
		HasMore:    query.GetHasMore(totalCount),
		Comments:   commentsList,
	}, nil
}
