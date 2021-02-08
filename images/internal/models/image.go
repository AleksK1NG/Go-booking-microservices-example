package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	imageService "github.com/AleksK1NG/hotels-mocroservices/images-microservice/proto/image"
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

func (i *Image) ToProto() *imageService.Image {
	return &imageService.Image{
		ImageID:    i.ImageID.String(),
		ImageURL:   i.ImageURL,
		IsUploaded: i.IsUploaded,
		CreatedAt:  timestamppb.New(i.CreatedAt),
	}
}

// UpdateHotelImageMsg
type UpdateHotelImageMsg struct {
	HotelID uuid.UUID `json:"hotel_id"`
	Image   string    `json:"image,omitempty"`
}
