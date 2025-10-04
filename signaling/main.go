package main

import (
	"encoding/json"
	"github.com/deva-labs/univia-api/signaling/services"
	"github.com/deva-labs/univia-api/signaling/store"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "2112"
	}

	http.HandleFunc("/ws", wsHandler)
	log.Printf("signaling listening %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Upgrade configuration for WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins (adjust for production)
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket: %v", err)
		return
	}
	defer conn.Close()

	ip := getUserIP(r)
	userID := getUserIDFromRedis(ip) // call Redis, nếu chưa có thì nil

	if userID == "" {
		log.Printf("No user_id for ip: %s", ip)
		conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"unauthorized"}`))
		return
	}

	// bind socket với user_id
	store.SetUserSocket(userID, conn)
	defer store.RemoveUserSocket(userID)

	log.Printf("Client connected: ip=%s userID=%s", ip, userID)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading: %v", err)
			break
		}

		var wsMessage services.WebSocketMessage
		if err := json.Unmarshal(message, &wsMessage); err != nil {
			log.Printf("Invalid message: %v", err)
			continue
		}

		services.HandleMessage(conn, messageType, wsMessage)
	}
}

func getUserIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return ip
}

func getUserIDFromRedis(ip string) string {
	val, err := redisClient.Get(context.Background(), "user:"+ip).Result()
	if err != nil {
		return ""
	}
	return val
}
