package api

import (
	"github.com/google/uuid"
)

type OErrorPacket = Packet[string]

const ErrorPacketType OPacketType = "server_error"

// TODO: 400 vs 500 class errors

func SendErr(c ClientInterface, clientId uuid.UUID, err error) {
	errorMessage := mp(ErrorPacketType, err.Error())

	Send(c, clientId, &errorMessage)
}
