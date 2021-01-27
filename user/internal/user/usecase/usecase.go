package usecase

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/AleksK1NG/hotels-mocroservices/user/internal/models"
	"github.com/AleksK1NG/hotels-mocroservices/user/internal/user"
)

// UserUseCase
type UserUseCase struct {
	userPGRepo user.PGRepository
}

// NewUserUseCase
func NewUserUseCase(userPGRepo user.PGRepository) *UserUseCase {
	return &UserUseCase{userPGRepo: userPGRepo}
}

// Register
func (u *UserUseCase) Register(ctx context.Context, user *models.User) (*models.UserResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.Register")
	defer span.Finish()

	if err := user.PrepareCreate(); err != nil {
		return nil, errors.Wrap(err, "user.PrepareCreate")
	}

	created, err := u.userPGRepo.Create(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, "user.PrepareCreate")
	}

	return created, err
}
