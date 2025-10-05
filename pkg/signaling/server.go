package signaling

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/deva-labs/univia/internal/infrastructure/redis"
	"github.com/deva-labs/univia/internal/signaling/services"
	"github.com/deva-labs/univia/internal/signaling/store"
	"github.com/deva-labs/univia/pkg/types"
)

type Server struct {
	Port string
}

func NewServer(port string) *Server {
	if port == "" {
		port = "2112"
	}
	return &Server{Port: port}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Add domain allowlist in production
		return true
	},
}

func (s *Server) Start() error {
	http.HandleFunc("/ws", s.handleWebSocket)
	log.Printf("[Signaling] Listening on port %s", s.Port)
	return http.ListenAndServe(":"+s.Port, nil)
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		http.Error(w, "Could not upgrade connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	ip := s.getUserIP(r)
	userID := s.getUserIDFromRedis(ip)
	if userID == uuid.Nil {
		log.Printf("Unauthorized WebSocket attempt: ip=%s", ip)
		conn.WriteJSON(map[string]string{"error": "unauthorized"})
		return
	}

	store.SetUserSocket(userID, conn)
	defer store.RemoveUserSocket(userID)
	log.Printf("Client connected: ip=%s userID=%s", ip, userID.String())

	s.handleMessages(conn)
}

func (s *Server) handleMessages(conn *websocket.Conn) {
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

func (s *Server) getUserIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func (s *Server) getUserIDFromRedis(ip string) uuid.UUID {
	val, err := redis.GetJSON[uuid.UUID](redis.Redis, "user:"+ip)
	if err != nil {
		log.Printf("Redis get user_id failed: %v", err)
		return uuid.Nil
	}
	return *val
}
