package api

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/Bismyth/game-server/pkg/db"
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
	Message  []byte
	SessonId uuid.UUID
}

type Client interface {
	Send([]uuid.UUID, []byte)
	Close(uuid.UUID)
}

type HandlerInput struct {
	C         Client
	Session   *db.Session
	SessionId uuid.UUID
	Packet    Packet[json.RawMessage]
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

func Send[T any](c Client, sessionId uuid.UUID, packet *Packet[T]) {
	c.Send([]uuid.UUID{sessionId}, MarshalPacket(packet))
}

func SendMany[T any](c Client, sessionIds []uuid.UUID, packet *Packet[T]) {
	c.Send(sessionIds, MarshalPacket(packet))
}

func HandleIncomingMessage(c Client, m *IRawMessage) {
	iPacket := Packet[json.RawMessage]{}

	var returnErr error

	returnErr = json.Unmarshal(m.Message, &iPacket)
	if returnErr != nil {
		SendErr(c, m.SessonId, returnErr)
		return
	}

	fn, ok := router[IPacketType(iPacket.Type)]
	if !ok {
		SendErr(c, m.SessonId, fmt.Errorf("unrecognized packet type"))
		return
	}

	s, err := db.GetSessionDetails(m.SessonId)
	if err != nil {
		SendErr(c, m.SessonId, err)
		return
	}

	returnErr = fn(HandlerInput{
		C:         c,
		Session:   s,
		SessionId: m.SessonId,
		Packet:    iPacket,
	})
	if returnErr != nil {
		SendErr(c, m.SessonId, returnErr)
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
