package models

import (
	uuid "github.com/satori/go.uuid"

	sessionService "github.com/AleksK1NG/hotels-mocroservices/api-gateway/proto/session"
)

// Session
type Session struct {
	UserID    uuid.UUID `json:"user_id"`
	SessionID string    `json:"session_id"`
}

// FromProto
func (s *Session) FromProto(session *sessionService.Session) (*Session, error) {
	userUUID, err := uuid.FromString(session.GetUserID())
	if err != nil {
		return nil, err
	}
	s.UserID = userUUID
	s.SessionID = session.GetSessionID()
	return s, nil
}
