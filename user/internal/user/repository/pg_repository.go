package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/user/internal/models"
	"github.com/AleksK1NG/hotels-mocroservices/user/pkg/utils"
)

type userPGRepository struct {
	db *pgxpool.Pool
}

func NewUserPGRepository(db *pgxpool.Pool) *userPGRepository {
	return &userPGRepository{db: db}
}

// Create new user
func (u *userPGRepository) Create(ctx context.Context, user *models.User) (*models.UserResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userPGRepository.Create")
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
func (u *userPGRepository) GetByID(ctx context.Context, userID uuid.UUID) (*models.UserResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userPGRepository.GetByID")
	defer span.Finish()

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

// GetByEmail
func (u *userPGRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userPGRepository.GetByEmail")
	defer span.Finish()

	var res models.User
	if err := u.db.QueryRow(ctx, getUserByEmail, email).Scan(
		&res.UserID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Password,
		&res.Avatar,
		&res.Role,
		&res.UpdatedAt,
		&res.CreatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}

	return &res, nil
}

// Update
func (u *userPGRepository) Update(ctx context.Context, user *models.UserUpdate) (*models.UserResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userPGRepository.Update")
	defer span.Finish()

	var res models.UserResponse
	if err := u.db.QueryRow(ctx, updateUserQuery, &user.FirstName, &user.LastName, &user.Email, &user.Role, &user.UserID).
		Scan(
			&res.UserID,
			&res.FirstName,
			&res.LastName,
			&res.Email,
			&res.Role,
			&res.Avatar,
			&res.UpdatedAt,
			&res.CreatedAt,
		); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}

	return &res, nil
}

func (u *userPGRepository) UpdateAvatar(ctx context.Context, msg models.UploadedImageMsg) (*models.UserResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userPGRepository.UpdateUploadedAvatar")
	defer span.Finish()

	log.Printf("REPO  IMAGE: %v", msg)
	var res models.UserResponse
	if err := u.db.QueryRow(ctx, updateAvatarQuery, &msg.ImageURL, &msg.UserID).Scan(
		&res.UserID,
		&res.FirstName,
		&res.LastName,
		&res.Email,
		&res.Role,
		&res.Avatar,
		&res.UpdatedAt,
		&res.CreatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}

	return &res, nil
}

func (u *userPGRepository) GetUsersByIDs(ctx context.Context, userIDs []string) ([]*models.UserResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userPGRepository.GetUsersByIDs")
	defer span.Finish()

	placeholders := utils.CreateSQLPlaceholders(len(userIDs))
	query := fmt.Sprintf("SELECT user_id, first_name, last_name, email, avatar, role, updated_at, created_at FROM users WHERE user_id IN (%v)", placeholders)

	args := make([]interface{}, len(userIDs))
	for i, id := range userIDs {
		args[i] = id
	}

	rows, err := u.db.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "db.Query")
	}
	defer rows.Close()

	users := make([]*models.UserResponse, 0, len(userIDs))
	for rows.Next() {
		var res models.UserResponse
		if err := rows.Scan(
			&res.UserID,
			&res.FirstName,
			&res.LastName,
			&res.Email,
			&res.Avatar,
			&res.Role,
			&res.UpdatedAt,
			&res.CreatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "db.Query")
		}
		users = append(users, &res)
	}

	return users, nil
}
