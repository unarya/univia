package types

type SessionMetadata struct {
	IP        string
	UserAgent string
}

type ResponseSession struct {
	AccessToken  string
	RefreshToken string
	SessionID    string
}
