package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/AleksK1NG/hotels-mocroservices/comments/proto/comments"
	userService "github.com/AleksK1NG/hotels-mocroservices/comments/proto/user"
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

// All Comments response with pagination
type CommentsList struct {
	TotalCount int        `json:"totalCount"`
	TotalPages int        `json:"totalPages"`
	Page       int        `json:"page"`
	Size       int        `json:"size"`
	HasMore    bool       `json:"hasMore"`
	Comments   []*Comment `json:"comments"`
}

// All Comments response with pagination
type CommentsFullList struct {
	TotalCount int                            `json:"totalCount"`
	TotalPages int                            `json:"totalPages"`
	Page       int                            `json:"page"`
	Size       int                            `json:"size"`
	HasMore    bool                           `json:"hasMore"`
	Comments   []*commentsService.CommentFull `json:"comments"`
}

// FullCommentsList
type FullCommentsList struct {
	CommentsList CommentsList
	UsersList    []UserResponse
}

// ToProto
func (h *CommentsList) ToProto() []*commentsService.Comment {
	commentsList := make([]*commentsService.Comment, 0, len(h.Comments))
	for _, hotel := range h.Comments {
		commentsList = append(commentsList, hotel.ToProto())
	}
	return commentsList
}

// ToHotelByIDProto
func (h *CommentsList) ToHotelByIDProto(users []*userService.User) []*commentsService.CommentFull {
	userMap := make(map[string]*userService.User, len(users))
	for _, user := range users {
		userMap[user.UserID] = user
	}

	commentsList := make([]*commentsService.CommentFull, 0, len(h.Comments))
	for _, comm := range h.Comments {
		user := userMap[comm.UserID.String()]

		commentsList = append(commentsList, &commentsService.CommentFull{
			CommentID: comm.CommentID.String(),
			HotelID:   comm.HotelID.String(),
			User: &commentsService.User{
				UserID:    user.UserID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
				Avatar:    user.Avatar,
				Role:      user.Role,
			},
			Message:   comm.Message,
			Photos:    comm.Photos,
			Rating:    comm.Rating,
			CreatedAt: timestamppb.New(*comm.CreatedAt),
			UpdatedAt: timestamppb.New(*comm.UpdatedAt),
		})
	}
	return commentsList
}
