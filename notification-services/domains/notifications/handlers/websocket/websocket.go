package websocket

import (
	"errors"
	"log"
	"net/http"
	"notification-services/config"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, you should have a whitelist of allowed origins.
		return true
	},
}

// Client represents a connected user.
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	userID uuid.UUID
}

// Hub maintains the set of active clients.
type Hub struct {
	clients    map[uuid.UUID]*websocket.Conn
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]*websocket.Conn),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's event loop.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client.userID] = client.conn
			h.mutex.Unlock()
			log.Printf("Client registered: %s", client.userID)

		case client := <-h.unregister:
			h.mutex.Lock()
			if currentConn, ok := h.clients[client.userID]; ok && currentConn == client.conn {
				delete(h.clients, client.userID)
			}
			h.mutex.Unlock()
			log.Printf("Client unregistered: %s", client.userID)
		}
	}
}

// PushToUser sends a message to a specific user if they are connected.
func (h *Hub) PushToUser(userID uuid.UUID, message interface{}) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if conn, ok := h.clients[userID]; ok {
		if err := conn.WriteJSON(message); err != nil {
			log.Printf("Error pushing to user %s: %v", userID, err)
		}
	}
}

// ValidateJWT validates the token and extracts claims. This should be in a shared package eventually.
func ValidateJWT(tokenString, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// serveWs handles the HTTP upgrade to a WebSocket connection.
func ServeWs(hub *Hub, c *gin.Context, cfg *config.Config) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token query parameter is missing"})
		return
	}

	claims, err := ValidateJWT(token, cfg.Server.JWTSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	userIDstr, ok := claims["user_id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id claim in token"})
		return
	}

	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Malformed user ID in token"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection for user %s: %v", userID, err)
		return
	}

	client := &Client{hub: hub, conn: conn, userID: userID}
	hub.register <- client

	// This goroutine listens for the connection to close.
	go func() {
		defer func() {
			hub.unregister <- client
			conn.Close()
		}()
		for {
			// Read messages from the client to detect a closed connection.
			if _, _, err := conn.NextReader(); err != nil {
				break
			}
		}
	}()
}
