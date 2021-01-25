package models

import uuid "github.com/satori/go.uuid"

// Session model
type Session struct {
	SessionID string    `json:"session_id"`
	UserID    uuid.UUID `json:"user_id"`
}
