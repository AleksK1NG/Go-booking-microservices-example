package models

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	uuid "github.com/satori/go.uuid"

	commentsService "github.com/AleksK1NG/hotels-mocroservices/api-gateway/proto/comments"
	userService "github.com/AleksK1NG/hotels-mocroservices/api-gateway/proto/user"
)

// User
type UserResponse struct {
	UserID    uuid.UUID  `json:"user_id"`
	FirstName string     `json:"first_name" validate:"required,min=3,max=25"`
	LastName  string     `json:"last_name" validate:"required,min=3,max=25"`
	Email     string     `json:"email" validate:"required,email"`
	Role      *string    `json:"role"`
	Avatar    *string    `json:"avatar" validate:"max=250" swaggertype:"string"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// UserFromProtoRes
func UserFromProtoRes(user *userService.User) (*UserResponse, error) {
	userUUID, err := uuid.FromString(user.GetUserID())
	if err != nil {
		return nil, err
	}

	createdAt, err := ptypes.Timestamp(user.CreatedAt)
	if err != nil {
		return nil, err
	}

	updatedAt, err := ptypes.Timestamp(user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		UserID:    userUUID,
		FirstName: user.GetFirstName(),
		LastName:  user.GetLastName(),
		Email:     user.GetEmail(),
		Role:      &user.Role,
		Avatar:    &user.Avatar,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}, nil
}

// UserProtoRes
type CommentUser struct {
	UserID    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	Role      string `json:"role"`
}

// FromProto
func (u *CommentUser) FromProto(user *commentsService.User) {
	u.UserID = user.UserID
	u.FirstName = user.FirstName
	u.LastName = user.LastName
	u.Email = user.Email
	u.Avatar = user.Avatar
	u.Role = user.Role
}
