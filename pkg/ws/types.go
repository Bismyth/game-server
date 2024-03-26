package ws

import (
	"encoding/json"

	"github.com/google/uuid"
)

type IncomingPacket struct {
	Message IncomingMessage
	Client  *Client
}

type IncomingMessage struct {
	Type MessageType     `json:"type"`
	Data json.RawMessage `json:"data"`
}

type OutgoingMessage struct {
	Type MessageType     `json:"type"`
	Data json.RawMessage `json:"data"`
}

type ChatMessage struct {
	Message string `json:"message"`
	Sender  string `json:"sender"`
}

type InitMessage struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type MessageType string

const (
	NamePacketType MessageType = "name"
)

const IdMessageType MessageType = "id"

type IdMessagePayload uuid.UUID
