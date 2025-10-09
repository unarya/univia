package types

import "github.com/google/uuid"

type SessionMetadata struct {
	IP        string
	UserAgent string
}

type ResponseSession struct {
	AccessToken  string
	RefreshToken string
	SessionID    uuid.UUID
	UserID       uuid.UUID
}
