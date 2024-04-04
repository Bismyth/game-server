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

func SendErr(c Client, userId uuid.UUID, err error) {
	errorMessage := CreateErrorPacket(err)

	Send(c, userId, &errorMessage)
}

const pt_ErrorGame = "server_error_game"

func SendGameErr(c Client, userId uuid.UUID, err error) {
	packet := mp(pt_ErrorGame, err.Error())

	Send(c, userId, &packet)
}
