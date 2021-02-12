package models

import commentsService "github.com/AleksK1NG/hotels-mocroservices/api-gateway/proto/comments"

// UserProtoRes
type CommentUser struct {
	UserID    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	Role      string `json:"role"`
}

// FromProto
func (u CommentUser) FromProto(user *commentsService.User) {
	u.UserID = user.UserID
	u.FirstName = user.FirstName
	u.LastName = user.LastName
	u.Email = user.Email
	u.Avatar = user.Avatar
	u.Role = user.Role
}
