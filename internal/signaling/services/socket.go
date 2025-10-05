package services

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/deva-labs/univia/internal/signaling/store"
	"github.com/deva-labs/univia/pkg/types"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// WebSocketMessage represents the JSON message format

// HandleMessage processes incoming WebSocket messages
func HandleMessage(conn *websocket.Conn, messageType int, wsMessage types.WebSocketMessage) error {
	var response types.WebSocketMessage

	switch wsMessage.Type {
	case "notice":
		response = types.WebSocketMessage{
			Type:    "notice",
			Message: "Received: " + wsMessage.Message,
		}
	case "bye":
		response = types.WebSocketMessage{
			Type:    "bye",
			Message: "Goodbye: " + wsMessage.Message,
		}
	case "ping":
		response = types.WebSocketMessage{
			Type:    "ping",
			Message: "Pong: " + wsMessage.Message,
		}
	default:
		response = types.WebSocketMessage{
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
func sendJSONMessage(conn *websocket.Conn, messageType int, response types.WebSocketMessage) error {
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return err
	}
	return conn.WriteMessage(messageType, responseJSON)
}

// SendMessageToUser is a function using socket to send message
func SendMessageToUser(userID uuid.UUID, message types.WebSocketMessage) error {
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
