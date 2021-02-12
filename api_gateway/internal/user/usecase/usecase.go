package usecase

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/internal/models"
	"github.com/AleksK1NG/hotels-mocroservices/api-gateway/pkg/logger"
	sessionService "github.com/AleksK1NG/hotels-mocroservices/api-gateway/proto/session"
	userService "github.com/AleksK1NG/hotels-mocroservices/api-gateway/proto/user"
)

// UserUseCase
type userUseCase struct {
	sessClient  sessionService.AuthorizationServiceClient
	userService userService.UserServiceClient
	logger      logger.Logger
}

// NewUserUseCase
func NewUserUseCase(
	sessClient sessionService.AuthorizationServiceClient,
	userService userService.UserServiceClient,
	logger logger.Logger,
) *userUseCase {
	return &userUseCase{sessClient: sessClient, userService: userService, logger: logger}
}

// GetByID
func (u *userUseCase) GetByID(ctx context.Context, userUUID uuid.UUID) (*models.UserResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userUseCase.GetByID")
	defer span.Finish()

	user, err := u.userService.GetUserByID(ctx, &userService.GetByIDRequest{UserID: userUUID.String()})
	if err != nil {
		return nil, errors.Wrap(err, "userService.GetUserByID")
	}

	res, err := models.UserFromProtoRes(user.GetUser())
	if err != nil {
		return nil, errors.Wrap(err, "UserFromProtoRes")
	}

	return res, nil
}

// GetSessionByID
func (u *userUseCase) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "userUseCase.GetSessionByID")
	defer span.Finish()

	sessionByID, err := u.sessClient.GetSessionByID(ctx, &sessionService.GetSessionByIDRequest{SessionID: sessionID})
	if err != nil {
		return nil, errors.Wrap(err, "sessClient.GetSessionByID")
	}

	sess := &models.Session{}
	sess, err = sess.FromProto(sessionByID.GetSession())
	if err != nil {
		return nil, errors.Wrap(err, "sess.FromProto")
	}

	return sess, nil
}
