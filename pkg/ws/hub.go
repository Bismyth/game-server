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
	clients map[*Client]bool

	// Client UUID map
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

type ClientInterface struct {
	clientIds map[uuid.UUID]*Client
}

func (c *ClientInterface) Send(ids []uuid.UUID, data []byte) {
	for _, id := range ids {
		client, ok := c.clientIds[id]
		if !ok {
			log.Printf("trying to send packet to unknown id: %v", id)
			continue
		}
		client.send <- data
	}
}

func newClientInterface() *ClientInterface {
	nci := ClientInterface{
		clientIds: make(map[uuid.UUID]*Client),
	}

	return &nci
}

func (h *Hub) Run() {

	clientInterface := newClientInterface()

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			clientInterface.clientIds[client.id] = client
			// TODO: Check for duplicate connection
			go api.SendUserInit(clientInterface, client.id)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.clientIds, client.id)
				close(client.send)
			}
		case message := <-h.broadcast:
			go api.HandleIncomingMessage(clientInterface, message)
		}
	}
}
