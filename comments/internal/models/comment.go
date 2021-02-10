package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	commentsService "github.com/AleksK1NG/hotels-mocroservices/comments/proto"
)

// Comment
type Comment struct {
	CommentID uuid.UUID  `json:"comment_id"`
	HotelID   uuid.UUID  `json:"hotel_id"`
	UserID    uuid.UUID  `json:"user_id"`
	Message   string     `json:"message" validate:"required,min=5,max=500"`
	Photos    []string   `json:"photos,omitempty"`
	Rating    float64    `json:"rating" validate:"required,min=0,max=10"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// ToProto
func (c *Comment) ToProto() *commentsService.Comment {
	return &commentsService.Comment{
		CommentID: c.CommentID.String(),
		HotelID:   c.HotelID.String(),
		UserID:    c.UserID.String(),
		Message:   c.Message,
		Photos:    c.Photos,
		Rating:    c.Rating,
		CreatedAt: timestamppb.New(*c.CreatedAt),
		UpdatedAt: timestamppb.New(*c.UpdatedAt),
	}
}
