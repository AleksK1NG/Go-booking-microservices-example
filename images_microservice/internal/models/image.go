package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Image model
type Image struct {
	ImageID    uuid.UUID `json:"image_id"`
	ImageURL   string    `json:"image_url"`
	IsUploaded bool      `json:"is_uploaded"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Event message for upload image
type UploadImageMsg struct {
	ImageID    uuid.UUID `json:"image_id"`
	UserID     uuid.UUID `json:"user_id"`
	ImageURL   string    `json:"image_url"`
	IsUploaded bool      `json:"is_uploaded"`
}

// Event message for create image
type CreateImageMsg struct {
	ImageURL   string `json:"image_url"`
	IsUploaded bool   `json:"is_uploaded"`
}
