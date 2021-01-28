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
		return nil, errors.Wrap(err, "userPGRepo.Create")
	}

	return created, err
}

// Login
func (u *UserUseCase) Login(ctx context.Context, login models.Login) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.Login")
	defer span.Finish()

	userByEmail, err := u.userPGRepo.GetByEmail(ctx, login.Email)
	if err != nil {
		return nil, errors.Wrap(err, "userPGRepo.GetByEmail")
	}

	if err := userByEmail.ComparePasswords(login.Password); err != nil {
		return nil, errors.Wrap(err, "UserUseCase.ComparePasswords")
	}

	userByEmail.SanitizePassword()

	return userByEmail, nil
}

// CreateSession
func (u *UserUseCase) CreateSession(ctx context.Context, userID uuid.UUID) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserUseCase.CreateSession")
	defer span.Finish()

	session, err := u.sessClient.CreateSession(ctx, &sessionService.CreateSessionRequest{UserID: userID.String()})
	if err != nil {
		return "", errors.Wrap(err, "sessClient.CreateSession")
	}

	return session.GetSession().GetSessionID(), err
}
