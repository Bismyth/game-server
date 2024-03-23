package ws

import (
	"encoding/json"

	"github.com/google/uuid"
)

type BaseMessage struct {
	Type   string          `json:"type"`
	Data   json.RawMessage `json:"data"`
	client *Client         `json:"-"`
}

type ChatMessage struct {
	Message string `json:"message"`
	Sender  string `json:"sender"`
}

type InitMessage struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
