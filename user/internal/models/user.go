package models

import (
	"fmt"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/AleksK1NG/hotels-mocroservices/user/pkg/types"
	userService "github.com/AleksK1NG/hotels-mocroservices/user/proto/user"
)

// User
type User struct {
	UserID    uuid.UUID            `json:"user_id"`
	FirstName string               `json:"first_name" validate:"required,min=3,max=25"`
	LastName  string               `json:"last_name" validate:"required,min=3,max=25"`
	Email     string               `json:"email" validate:"required,email"`
	Password  string               `json:"password" validate:"required,min=6,max=250"`
	Avatar    types.NullJSONString `json:"avatar" validate:"max=250"`
	Role      *Role                `json:"role"`
	CreatedAt *time.Time           `json:"created_at"`
	UpdatedAt *time.Time           `json:"updated_at"`
}

// User
type UserResponse struct {
	UserID    uuid.UUID            `json:"user_id"`
	FirstName string               `json:"first_name" validate:"required,min=3,max=25"`
	LastName  string               `json:"last_name" validate:"required,min=3,max=25"`
	Email     string               `json:"email" validate:"required,email"`
	Role      *Role                `json:"role"`
	Avatar    types.NullJSONString `json:"avatar" validate:"max=250"`
	CreatedAt *time.Time           `json:"created_at"`
	UpdatedAt *time.Time           `json:"updated_at"`
}

type Role string

const (
	RoleGuest  Role = "guest"
	RoleMember Role = "member"
	RoleAdmin  Role = "admin"
)

func (e *Role) ToString() string {
	return string(*e)
}

func (e *Role) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Role(s)
	case string:
		*e = Role(s)
	default:
		return fmt.Errorf("unsupported scan type for Role: %T", src)
	}
	return nil
}

// Hash user password with bcrypt
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Compare user password and payload
func (u *User) ComparePasswords(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

// Sanitize user password
func (u *User) SanitizePassword() {
	u.Password = ""
}

// Prepare user for register
func (u *User) PrepareCreate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)

	if err := u.HashPassword(); err != nil {
		return err
	}
	return nil
}

func (r *UserResponse) ToProto() *userService.User {
	return &userService.User{
		UserID:    r.UserID.String(),
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Email,
		Avatar:    r.Avatar.String,
		Role:      r.Role.ToString(),
		CreatedAt: timestamppb.New(*r.CreatedAt),
		UpdatedAt: timestamppb.New(*r.UpdatedAt),
	}
}
