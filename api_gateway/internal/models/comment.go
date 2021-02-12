package models

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	uuid "github.com/satori/go.uuid"

	commentsService "github.com/AleksK1NG/hotels-mocroservices/api-gateway/proto/comments"
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

// All Comments response with pagination
type CommentsList struct {
	TotalCount int64          `json:"totalCount"`
	TotalPages int64          `json:"totalPages"`
	Page       int64          `json:"page"`
	Size       int64          `json:"size"`
	HasMore    bool           `json:"hasMore"`
	Comments   []*CommentFull `json:"comments"`
}

// CommentFromProto
func CommentFromProto(comment *commentsService.Comment) (*Comment, error) {
	commUUID, err := uuid.FromString(comment.CommentID)
	if err != nil {
		return nil, err
	}
	userUUID, err := uuid.FromString(comment.UserID)
	if err != nil {
		return nil, err
	}
	hotelUUID, err := uuid.FromString(comment.HotelID)
	if err != nil {
		return nil, err
	}

	createdAt, err := ptypes.Timestamp(comment.CreatedAt)
	if err != nil {
		return nil, err
	}

	updatedAt, err := ptypes.Timestamp(comment.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &Comment{
		CommentID: commUUID,
		HotelID:   hotelUUID,
		UserID:    userUUID,
		Message:   comment.Message,
		Photos:    comment.Photos,
		Rating:    comment.Rating,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}, nil
}

// CommentFull
type CommentFull struct {
	CommentID uuid.UUID    `json:"comment_id"`
	HotelID   uuid.UUID    `json:"hotel_id"`
	User      *CommentUser `json:"user"`
	Message   string       `json:"message"`
	Photos    []string     `json:"photos"`
	Rating    float64      `json:"rating"`
	CreatedAt *time.Time   `json:"createdAt"`
	UpdatedAt *time.Time   `json:"updatedAt"`
}

// CommentFullFromProto
func CommentFullFromProto(comm *commentsService.CommentFull) (*CommentFull, error) {
	commUUID, err := uuid.FromString(comm.GetCommentID())
	if err != nil {
		return nil, err
	}

	hotelUUID, err := uuid.FromString(comm.GetHotelID())
	if err != nil {
		return nil, err
	}

	createdAt, err := ptypes.Timestamp(comm.CreatedAt)
	if err != nil {
		return nil, err
	}

	updatedAt, err := ptypes.Timestamp(comm.UpdatedAt)
	if err != nil {
		return nil, err
	}

	user := &CommentUser{}
	user.FromProto(comm.GetUser())

	return &CommentFull{
		CommentID: commUUID,
		HotelID:   hotelUUID,
		User:      user,
		Message:   comm.GetMessage(),
		Photos:    comm.GetPhotos(),
		Rating:    comm.GetRating(),
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}, nil
}
