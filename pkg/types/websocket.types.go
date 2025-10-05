package types

type WebSocketMessage struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	ReceiverID string `json:"receiverId"`
}
