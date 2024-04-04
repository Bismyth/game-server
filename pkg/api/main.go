package api

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

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

type Client interface {
	Send([]uuid.UUID, []byte)
}

type HandlerInput struct {
	C      Client
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

func Send[T any](c Client, clientId uuid.UUID, packet *Packet[T]) {
	c.Send([]uuid.UUID{clientId}, MarshalPacket(packet))
}

func SendMany[T any](c Client, clientIds []uuid.UUID, packet *Packet[T]) {
	c.Send(clientIds, MarshalPacket(packet))
}

func HandleIncomingMessage(c Client, m *IRawMessage) {
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
			customStringParsing(),
		),
		Result: &output,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}

func customStringParsing() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}

		switch t {
		case reflect.TypeOf(uuid.UUID{}):
			return uuid.Parse(data.(string))
		case reflect.TypeOf(true):
			return data.(string) == "true", nil
		default:
			return data, nil
		}

	}
}
