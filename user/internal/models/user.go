package models

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

// User
type User struct {
	UserID    uuid.UUID  `json:"user_id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	Role      *Role      `json:"role"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type Role string

const (
	RoleGuest  Role = "guest"
	RoleMember Role = "member"
	RoleAdmin  Role = "admin"
)

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
