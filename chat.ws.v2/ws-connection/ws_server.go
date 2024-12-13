package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Server struct {
	Connections map[*websocket.Conn]bool
	mu          sync.Mutex
}

func NewServer() *Server {
	return &Server{
		Connections: make(map[*websocket.Conn]bool),
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins; modify if needed for security
		return true
	},
}

func (s *Server) HandleWSConnection(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer ws.Close()

	// Add connection to the map
	s.mu.Lock()
	s.Connections[ws] = true
	s.mu.Unlock()
	log.Println("New WebSocket connection established")

	// Remove connection from the map on disconnect
	defer func() {
		s.mu.Lock()
		delete(s.Connections, ws)
		s.mu.Unlock()
		log.Println("WebSocket connection closed")
	}()

	// Read messages from the WebSocket client
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Println("Client disconnected normally")
				break
			}
			log.Println("Error reading message:", err)
			continue
		}
		log.Printf("Received message: %s\n", msg)

		// Echo back a response
		if err := ws.WriteMessage(websocket.TextMessage, []byte("Message received successfully and is being processed.")); err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}

func main() {
	server := NewServer()

	http.HandleFunc("/ws", server.HandleWSConnection)

	log.Println("Starting WebSocket server on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}
}
