package services

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deva-labs/univia/signaling/store"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// WebSocketMessage represents the JSON message format
type WebSocketMessage struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	ReceiverID string `json:"receiverId"`
}

// HandleMessage processes incoming WebSocket messages
func HandleMessage(conn *websocket.Conn, messageType int, wsMessage WebSocketMessage) error {
	var response WebSocketMessage

	switch wsMessage.Type {
	case "notice":
		response = WebSocketMessage{
			Type:    "notice",
			Message: "Received: " + wsMessage.Message,
		}
	case "bye":
		response = WebSocketMessage{
			Type:    "bye",
			Message: "Goodbye: " + wsMessage.Message,
		}
	case "ping":
		response = WebSocketMessage{
			Type:    "ping",
			Message: "Pong: " + wsMessage.Message,
		}
	default:
		response = WebSocketMessage{
			Type:    "echo",
			Message: "Echo: " + wsMessage.Message,
		}
	}

	// Send the response as JSON
	if err := sendJSONMessage(conn, messageType, response); err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}
	return nil
}

// sendJSONMessage sends a JSON-encoded message to the client
func sendJSONMessage(conn *websocket.Conn, messageType int, response WebSocketMessage) error {
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return err
	}
	return conn.WriteMessage(messageType, responseJSON)
}

// SendMessageToUser is a function using socket to send message
func SendMessageToUser(userID uuid.UUID, message WebSocketMessage) error {
	conn, exists := store.GetUserSocket(userID)
	if !exists {
		return fmt.Errorf("user %s not connected", userID)
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, data)
}
