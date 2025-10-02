package store

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	UserSocketMap = make(map[uuid.UUID]*websocket.Conn)
	MapMutex      = sync.RWMutex{} // processing concurrent map access
)

func SetUserSocket(userID uuid.UUID, conn *websocket.Conn) {
	MapMutex.Lock()
	defer MapMutex.Unlock()
	UserSocketMap[userID] = conn
}

func GetUserSocket(userID uuid.UUID) (*websocket.Conn, bool) {
	MapMutex.RLock()
	defer MapMutex.RUnlock()
	conn, exists := UserSocketMap[userID]
	return conn, exists
}

func RemoveUserSocket(userID uuid.UUID) {
	MapMutex.Lock()
	defer MapMutex.Unlock()
	delete(UserSocketMap, userID)
}
