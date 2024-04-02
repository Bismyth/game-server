package interfaces

import (
	"github.com/google/uuid"
)

type Client interface {
	Send([]uuid.UUID, []byte)
}

type GameCommunication interface {
	SendEvent(data any)
	SendGlobal(data any)
	SendPlayer(playerId uuid.UUID, data any)
	ActionPrompt(playerId uuid.UUID, data any)
}
