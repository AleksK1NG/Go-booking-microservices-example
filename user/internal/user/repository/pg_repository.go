package repository

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/user/internal/models"
)

// UserPGRepository
type UserPGRepository struct {
	db *pgxpool.Pool
}

// NewUserPGRepository
func NewUserPGRepository(db *pgxpool.Pool) *UserPGRepository {
	return &UserPGRepository{db: db}
}

// Create new user
func (u *UserPGRepository) Create(ctx context.Context, user *models.User) (*models.UserResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserPGRepository.Create")
	defer span.Finish()

	var created models.UserResponse
	if err := u.db.QueryRow(
		ctx,
		createUserQuery,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Avatar,
		&user.Role,
	).Scan(&created.UserID, &created.FirstName, &created.LastName, &created.Email,
		&created.Avatar, &created.Role, &created.UpdatedAt, &created.CreatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}

	return &created, nil
}

// Get user by id
func (u *UserPGRepository) GetByID(ctx context.Context, userID uuid.UUID) (*models.UserResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserPGRepository.GetByID")
	defer span.Finish()

	getUserByIDQuery := `SELECT  user_id, first_name, last_name, email, avatar, role, updated_at, created_at  FROM users WHERE user_id = $1`
	var res models.UserResponse
	if err := u.db.QueryRow(ctx, getUserByIDQuery, userID).Scan(
		&res.UserID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Avatar,
		&res.Role,
		&res.UpdatedAt,
		&res.CreatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}

	return &res, nil
}
