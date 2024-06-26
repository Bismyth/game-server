package ws

import (
	"log"

	"github.com/Bismyth/game-server/pkg/api"
	"github.com/google/uuid"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients   map[*Client]bool
	clientIds map[uuid.UUID]*Client

	// Inbound messages from the clients.
	broadcast chan *api.IRawMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *api.IRawMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		clientIds:  make(map[uuid.UUID]*Client),
	}
}

func (h *Hub) Send(ids []uuid.UUID, data []byte) {
	for _, id := range ids {
		client, ok := h.clientIds[id]
		if !ok {
			log.Printf("trying to send packet to unknown id: %v", id)
			continue
		}
		client.send <- data
	}
}

var roomLeave = CloseMessage{
	Code:    3001,
	Message: "leave",
}

func (h *Hub) Close(id uuid.UUID) {
	client, ok := h.clientIds[id]
	if !ok {
		log.Printf("tried to close unknown session: %v", id)
	}
	client.leaveMessage = &roomLeave
	h.unregister <- client
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.clientIds[client.sessionId] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				err := api.HandleSessionClose(client.sessionId)
				if err != nil {
					log.Printf("failed to close session: %v", err)
				}
				delete(h.clients, client)
				delete(h.clientIds, client.sessionId)
				close(client.send)
			}
		case message := <-h.broadcast:
			go api.HandleIncomingMessage(h, message)
		}
	}
}
