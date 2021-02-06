package models

import (
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/AleksK1NG/hotels-mocroservices/hotels/pkg/types"
)

// Hotel model
type Hotel struct {
	HotelID       uuid.UUID         `json:"hotel_id"`
	Name          string            `json:"name"`
	Email         string            `json:"email,omitempty"`
	Country       string            `json:"country,omitempty"`
	City          string            `json:"city,omitempty"`
	Description   string            `json:"description,omitempty"`
	Location      string            `json:"location"`
	Rating        float64           `json:"rating"`
	Image         types.NullString  `json:"image"`
	Photos        []string          `json:"photos,omitempty"`
	CommentsCount int               `json:"comments_count,omitempty"`
	Latitude      types.NullFloat64 `json:"latitude,omitempty"`
	Longitude     types.NullFloat64 `json:"longitude,omitempty"`
	CreatedAt     *time.Time        `json:"created_at"`
	UpdatedAt     *time.Time        `json:"updated_at"`
}
