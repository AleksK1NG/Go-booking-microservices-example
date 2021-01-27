package grpc

import (
	"context"

	"github.com/opentracing/opentracing-go"
	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/user/internal/user"
	"github.com/AleksK1NG/hotels-mocroservices/user/pkg/grpc_errors"
	"github.com/AleksK1NG/hotels-mocroservices/user/pkg/logger"
	"github.com/AleksK1NG/hotels-mocroservices/user/proto/user"
)

// UserService
type UserService struct {
	userUC user.UseCase
	logger logger.Logger
}

// NewUserService
func NewUserService(userUC user.UseCase, logger logger.Logger) *UserService {
	return &UserService{userUC: userUC, logger: logger}
}

// GetUserByID
func (u *UserService) GetUserByID(ctx context.Context, r *userService.GetByIDRequest) (*userService.GetByIDResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserService.GetUserByID")
	defer span.Finish()

	userUUID, err := uuid.FromString(r.GetUserID())
	if err != nil {
		u.logger.Errorf("uuid.FromString: %v", err)
		return nil, grpc_errors.ErrorResponse(err, "uuid.FromString")
	}

	foundUser, err := u.userUC.GetByID(ctx, userUUID)
	if err != nil {
		u.logger.Errorf("uuid.FromString: %v", err)
		return nil, grpc_errors.ErrorResponse(err, "userUC.GetByID")
	}

	return &userService.GetByIDResponse{User: foundUser.ToProto()}, nil
}
