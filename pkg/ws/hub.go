package ws

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/google/uuid"
	"github.com/goombaio/namegenerator"
)

var nameGenerator = namegenerator.NewNameGenerator(time.Now().UTC().UnixNano())

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Client UUID map
	clientIds map[uuid.UUID]*Client
	// Inbound messages from the clients.
	broadcast chan *BaseMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *BaseMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		clientIds:  make(map[uuid.UUID]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.clientIds, client.id)
				close(client.send)
			}
		case message := <-h.broadcast:
			if message.Type == "id" {
				h.handleId(message)
				break
			}

			if message.client.id == uuid.Nil {
				break
			}

			switch message.Type {
			case "createLobby":
				h.createLobby(message)
			case "joinLobby":
				h.joinLobby(message)
			case "leaveLobby":
				h.leaveLobby(message)
			case "chat":
				h.handleChat(message)
			case "name":
				h.handleNameChange(message)
			default:
				fmt.Printf("unknown message type")
			}
		}
	}
}

func (h *Hub) handleChat(message *BaseMessage) error {
	var chatMessage ChatMessage

	err := json.Unmarshal(message.Data, &chatMessage)
	if err != nil {
		fmt.Println("couldn't unmarshal")
	}

	client, err := db.GetClient(message.client.id)
	if err != nil {
		return err
	}
	chatMessage.Sender = client.Name

	raw, err := json.Marshal(chatMessage)
	if err != nil {
		return err
	}
	message.Data = raw
	for client := range h.clients {
		select {
		case client.send <- message:

		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
	return nil
}

func (h *Hub) handleNameChange(message *BaseMessage) error {
	var name string
	err := json.Unmarshal(message.Data, &name)
	if err != nil {
		return err
	}

	client, err := db.GetClient(message.client.id)
	if err != nil {
		return err
	}
	client.Name = name
	client.Save()

	return nil
}

func (h *Hub) handleId(message *BaseMessage) error {
	var id uuid.UUID
	err := json.Unmarshal(message.Data, &id)
	if err != nil {
		return err
	}

	if client, ok := h.clientIds[id]; ok {
		close(client.send)
		delete(h.clients, client)
		delete(h.clientIds, id)
	}

	client, err := db.GetClient(id)

	if err != nil {
		newId, err := uuid.NewV7()
		if err != nil {
			return err
		}

		client = &db.User{Id: newId, Name: nameGenerator.Generate()}

	}

	client.Save()

	message.client.id = client.Id

	data, err := json.Marshal(client.Id.String())
	if err != nil {
		return err
	}

	websocketMessage := &BaseMessage{
		Type: "id",
		Data: data,
	}

	message.client.send <- websocketMessage

	h.clientIds[client.Id] = message.client

	return nil
}
