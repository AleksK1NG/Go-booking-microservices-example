package models

// CsrfToken
type CsrfToken struct {
	SessionID string `json:"session_id"`
	Timestamp int64  `json:"timestamp"`
}
