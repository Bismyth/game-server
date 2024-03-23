package ws

import (
	"encoding/json"
	"fmt"
	"log"
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
			err := h.handleMessageType(message)
			if err != nil {
				errorPayload, jErr := json.Marshal(err.Error())
				if jErr != nil {
					log.Printf("failed to marshal error to websocket: %v", jErr)
				}

				errorMessage := BaseMessage{
					Type: "error",
					Data: errorPayload,
				}

				message.client.send <- &errorMessage
			}
		}
	}
}

func (h *Hub) handleMessageType(message *BaseMessage) error {
	var err error

	if message.Type == "id" {
		err = h.handleId(message)
		return err
	}

	if message.client.id == uuid.Nil {
		return fmt.Errorf("no id in message")
	}

	switch message.Type {
	case "createLobby":
		err = h.createLobby(message)
	case "joinLobby":
		err = h.joinLobby(message)
	case "leaveLobby":
		err = h.leaveLobby(message)
	case "chat":
		err = h.handleChat(message)
	case "name":
		err = h.handleNameChange(message)
	default:
		err = fmt.Errorf("unknown message type")
	}

	return err
}

func (h *Hub) handleChat(message *BaseMessage) error {
	var chatMessage ChatMessage

	err := json.Unmarshal(message.Data, &chatMessage)
	if err != nil {
		return err
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

	nameBytes, err := json.Marshal(client.Name)
	if err != nil {
		return err
	}

	response := BaseMessage{
		Type: "name",
		Data: nameBytes,
	}

	message.client.send <- &response

	return nil
}

func (h *Hub) handleId(message *BaseMessage) error {
	var user InitMessage
	err := json.Unmarshal(message.Data, &user)
	if err != nil {
		return err
	}

	if client, ok := h.clientIds[user.Id]; ok {
		close(client.send)
		delete(h.clients, client)
		delete(h.clientIds, user.Id)
	}

	client, err := db.GetClient(user.Id)

	if err != nil {
		newId, err := uuid.NewV7()
		if err != nil {
			return err
		}

		client = &db.User{Id: newId, Name: nameGenerator.Generate()}

	}

	client.Save()

	message.client.id = client.Id

	initMessage := InitMessage{
		Id:   client.Id,
		Name: client.Name,
	}

	data, err := json.Marshal(&initMessage)
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
