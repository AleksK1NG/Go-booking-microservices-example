package usecase

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/user/internal/models"
	"github.com/AleksK1NG/hotels-mocroservices/user/internal/user"
	sessionService "github.com/AleksK1NG/hotels-mocroservices/user/proto/session"
)

// UserUseCase
type UserUseCase struct {
	userPGRepo user.PGRepository
	sessClient sessionService.AuthorizationServiceClient
}

// NewUserUseCase
func NewUserUseCase(userPGRepo user.PGRepository, sessClient sessionService.AuthorizationServiceClient) *UserUseCase {
	return &UserUseCase{userPGRepo: userPGRepo, sessClient: sessClient}
}

// GetByID
func (u *UserUseCase) GetByID(ctx context.Context, userID uuid.UUID) (*models.UserResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.GetByID")
	defer span.Finish()

	return u.userPGRepo.GetByID(ctx, userID)
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
