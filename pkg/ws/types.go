package ws

import "encoding/json"

type BaseMessage struct {
	Type   string          `json:"type"`
	Data   json.RawMessage `json:"data"`
	client *Client         `json:"-"`
}

type ChatMessage struct {
	Message string `json:"message"`
	Sender  string `json:"sender"`
}
