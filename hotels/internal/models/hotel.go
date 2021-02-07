package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	hotelsService "github.com/AleksK1NG/hotels-mocroservices/hotels/proto/hotels"
)

// Hotel model
type Hotel struct {
	HotelID       uuid.UUID  `json:"hotel_id"`
	Name          string     `json:"name" validate:"required,min=3,max=25"`
	Email         string     `json:"email,omitempty" validate:"required,email"`
	Country       string     `json:"country,omitempty" validate:"required,min=3,max=25"`
	City          string     `json:"city,omitempty" validate:"required,min=3,max=25"`
	Description   string     `json:"description,omitempty" validate:"required,min=10,max=250"`
	Location      string     `json:"location" validate:"required,min=10,max=250"`
	Rating        float64    `json:"rating" validate:"required,min=0,max=10"`
	Image         *string    `json:"image,omitempty"`
	Photos        []string   `json:"photos,omitempty"`
	CommentsCount int        `json:"comments_count,omitempty"`
	Latitude      *float64   `json:"latitude,omitempty"`
	Longitude     *float64   `json:"longitude,omitempty"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

func (h *Hotel) ToProto() *hotelsService.Hotel {
	return &hotelsService.Hotel{
		HotelID:       h.HotelID.String(),
		Name:          h.Name,
		Email:         h.Email,
		Country:       h.Country,
		City:          h.City,
		Description:   h.Description,
		Image:         *h.Image,
		Photos:        h.Photos,
		CommentsCount: int64(h.CommentsCount),
		Latitude:      *h.Latitude,
		Longitude:     *h.Longitude,
		Location:      h.Location,
		CreatedAt:     timestamppb.New(*h.CreatedAt),
		UpdatedAt:     timestamppb.New(*h.UpdatedAt),
	}
}
