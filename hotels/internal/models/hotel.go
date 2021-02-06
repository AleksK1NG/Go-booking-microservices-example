package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/AleksK1NG/hotels-mocroservices/hotels/pkg/types"
	hotelsService "github.com/AleksK1NG/hotels-mocroservices/hotels/proto/hotels"
)

// Hotel model
type Hotel struct {
	HotelID       uuid.UUID         `json:"hotel_id"`
	Name          string            `json:"name" validate:"required,min=3,max=25"`
	Email         string            `json:"email,omitempty" validate:"required,email"`
	Country       string            `json:"country,omitempty" validate:"required,min=3,max=25"`
	City          string            `json:"city,omitempty" validate:"required,min=3,max=25"`
	Description   string            `json:"description,omitempty" validate:"required,min=10,max=250"`
	Location      string            `json:"location" validate:"required,min=10,max=25"`
	Rating        float64           `json:"rating" validate:"required,min=0,max=10"`
	Image         types.NullString  `json:"image,omitempty"`
	Photos        []string          `json:"photos,omitempty"`
	CommentsCount int               `json:"comments_count,omitempty"`
	Latitude      types.NullFloat64 `json:"latitude,omitempty"`
	Longitude     types.NullFloat64 `json:"longitude,omitempty"`
	CreatedAt     *time.Time        `json:"created_at"`
	UpdatedAt     *time.Time        `json:"updated_at"`
}

func (h *Hotel) ToProto() *hotelsService.Hotel {
	return &hotelsService.Hotel{
		Name:          h.Name,
		Email:         h.Email,
		Country:       h.Country,
		City:          h.City,
		Description:   h.Description,
		Image:         h.Image.String,
		Photos:        h.Photos,
		CommentsCount: int64(h.CommentsCount),
		Latitude:      h.Latitude.Float64,
		Longitude:     h.Longitude.Float64,
		Location:      h.Location,
		CreatedAt:     timestamppb.New(*h.CreatedAt),
		UpdatedAt:     timestamppb.New(*h.UpdatedAt),
	}
}
