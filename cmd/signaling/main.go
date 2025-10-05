package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/deva-labs/univia/internal/infrastructure/redis"
	"github.com/deva-labs/univia/internal/signaling/services"
	"github.com/deva-labs/univia/internal/signaling/store"
	"github.com/deva-labs/univia/pkg/types"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func main() {
	port := os.Getenv("SIGNALING_PORT")
	if port == "" {
		port = "2112"
	}

	redis.ConnectRedis()
	http.HandleFunc("/ws", handleWebSocket)

	log.Printf("[Signaling] Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// ---- WebSocket ----
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Add domain-based allowlist for production
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		http.Error(w, "Could not upgrade connection", http.StatusBadRequest)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("Error closing websocket connection: %v", err)
		}
	}(conn)

	ip := getUserIP(r)
	userID := getUserIDFromRedis(ip)
	if userID == uuid.Nil {
		log.Printf("Unauthorized WebSocket attempt: ip=%s", ip)
		err := conn.WriteJSON(map[string]string{"error": "unauthorized"})
		if err != nil {
			return
		}
		return
	}

	// Store and manage socket
	store.SetUserSocket(userID, conn)
	defer store.RemoveUserSocket(userID)
	log.Printf("Client connected: ip=%s userID=%s", ip, userID.String())

	handleMessages(conn)
}

func handleMessages(conn *websocket.Conn) {
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		var wsMsg types.WebSocketMessage
		if err := json.Unmarshal(msg, &wsMsg); err != nil {
			log.Printf("Invalid message: %v", err)
			continue
		}

		if err := services.HandleMessage(conn, messageType, wsMsg); err != nil {
			log.Printf("HandleMessage error: %v", err)
		}
	}
}

// ---- Helpers ----

func getUserIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func getUserIDFromRedis(ip string) uuid.UUID {
	val, err := redis.GetJSON[uuid.UUID](redis.Redis, "user:"+ip)
	if err != nil {
		log.Printf("Redis get user_id failed: %v", err)
		return uuid.Nil
	}
	return *val
}
