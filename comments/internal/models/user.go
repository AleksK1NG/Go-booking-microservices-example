package models

// User
type UserResponse struct {
	UserID    string  `json:"user_id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email"`
	Avatar    *string `json:"avatar"`
	Role      *string `json:"role"`
}
