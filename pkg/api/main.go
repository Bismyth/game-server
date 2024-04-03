package api

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/Bismyth/game-server/pkg/interfaces"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
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

type HandlerInput struct {
	C      interfaces.Client
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

func MarshalPacket[T any](packet *Packet[T]) []byte {
	data, err := json.Marshal(packet)
	if err != nil {
		log.Printf("failed to marshal api packet: %v", err)
	}
	return data
}

func Send[T any](c interfaces.Client, clientId uuid.UUID, packet *Packet[T]) {
	c.Send([]uuid.UUID{clientId}, MarshalPacket(packet))
}

func SendMany[T any](c interfaces.Client, clientIds []uuid.UUID, packet *Packet[T]) {
	c.Send(clientIds, MarshalPacket(packet))
}

func HandleIncomingMessage(c interfaces.Client, m *IRawMessage) {
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

func decode(input, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			stringToUUIDHookFunc(),
		),
		Result: &output,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}

func stringToUUIDHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(uuid.UUID{}) {
			return data, nil
		}

		return uuid.Parse(data.(string))
	}
}
