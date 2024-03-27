package api

import (
	"github.com/google/uuid"
)

type OErrorPacket = Packet[string]

const ErrorPacketType OPacketType = "server_error"

// TODO: 400 vs 500 class errors

func CreateErrorPacket(err error) Packet[string] {
	return mp(ErrorPacketType, err.Error())
}

func SendErr(c ClientInterface, userId uuid.UUID, err error) {
	errorMessage := CreateErrorPacket(err)

	Send(c, userId, &errorMessage)
}
