package api

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
)

/*
- Modify hub to call this package and pass in a client communication interface™
- Calles to this package are made in a subprocess
- This pacakage defines message types, vaildators and handlers
*/

// Incoming Packet
type packetType string
type IPacketType packetType
type OPacketType packetType

type Packet[T interface{}] struct {
	Type packetType `json:"type"`
	Data T          `json:"data"`
}

type IRawMessage struct {
	Message []byte
	UserId  uuid.UUID
}

type ClientInterface interface {
	Send([]uuid.UUID, []byte)
}

type HandlerInput struct {
	C      ClientInterface
	UserId uuid.UUID
	Packet Packet[json.RawMessage]
}

func mp[T any](oPacketType OPacketType, data T) Packet[T] {
	return Packet[T]{
		Type: packetType(oPacketType),
		Data: data,
	}
}

func hp[T any](packet Packet[json.RawMessage]) (*T, error) {
	var decoded T
	err := json.Unmarshal(packet.Data, &decoded)
	if err != nil {
		return nil, err
	}
	return &decoded, nil
}

func MarsahlPacket[T any](packet *Packet[T]) []byte {
	data, err := json.Marshal(packet)
	if err != nil {
		log.Printf("failed to marshal api packet: %v", err)
	}
	return data
}

func Send[T any](c ClientInterface, clientId uuid.UUID, packet *Packet[T]) {
	c.Send([]uuid.UUID{clientId}, MarsahlPacket(packet))
}

func SendMany[T any](c ClientInterface, clientIds []uuid.UUID, packet *Packet[T]) {
	c.Send(clientIds, MarsahlPacket(packet))
}

func HandleIncomingMessage(c ClientInterface, m *IRawMessage) {
	iPacket := Packet[json.RawMessage]{}

	var returnErr error

	returnErr = json.Unmarshal(m.Message, &iPacket)
	if returnErr != nil {
		SendErr(c, m.UserId, returnErr)
		return
	}

	fn, ok := router[IPacketType(iPacket.Type)]
	if !ok {
		SendErr(c, m.UserId, fmt.Errorf("unrecognized packet type"))
		return
	}

	returnErr = fn(HandlerInput{
		C:      c,
		UserId: m.UserId,
		Packet: iPacket,
	})
	if returnErr != nil {
		SendErr(c, m.UserId, returnErr)
	}
}
