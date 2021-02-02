package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Event message for uploaded images
type UploadedImageMsg struct {
	ImageID    uuid.UUID `json:"image_id"`
	UserID     uuid.UUID `json:"user_id"`
	ImageURL   string    `json:"image_url"`
	IsUploaded bool      `json:"is_uploaded"`
	CreatedAt  time.Time `json:"created_at"`
}

// Event message for upload user avatar
type UpdateAvatarMsg struct {
	UserID      uuid.UUID `json:"user_id"`
	ContentType string    `json:"content_type"`
	Body        []byte
}

// Image model
type Image struct {
	ImageID    uuid.UUID `json:"image_id"`
	ImageURL   string    `json:"image_url"`
	IsUploaded bool      `json:"is_uploaded"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
