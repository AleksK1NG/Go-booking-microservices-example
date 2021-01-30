package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Image model
type Image struct {
	ImageID    uuid.UUID `json:"image_id"`
	ImageURL   string    `json:"image_url"`
	CreatedAt  time.Time `json:"created_at"`
	IsUploaded bool      `json:"is_uploaded"`
}
