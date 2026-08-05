package game

import (
	"encoding/json"
	"fmt"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/Bismyth/game-server/pkg/interfaces"
	"github.com/google/uuid"

	liarsdice "github.com/Bismyth/game-server/pkg/game/liars-dice"
	skull "github.com/Bismyth/game-server/pkg/game/skull"
)

type GameHandler interface {
	HandleAction(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID, data json.RawMessage) error
	HandleReady(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID) error
	New(gameId uuid.UUID, options []byte) error
	DefaultOptions() interface{}
	HandleLeave(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID) error
	Cleanup(gameId uuid.UUID) error
}

var gameHandlers map[string]GameHandler = map[string]GameHandler{
	liarsdice.Code: liarsdice.New(),
	skull.Code:     skull.New(),
}

type SharedState struct {
	Id uuid.UUID
}

func getGameType(gameId uuid.UUID) (GameHandler, error) {
	gameType, err := db.GetRoomProperty[string](gameId, "gameType")
	if err != nil {
		return nil, err
	}

	h, ok := gameHandlers[gameType]
	if !ok {
		return nil, fmt.Errorf("game type not found in list")
	}
	return h, nil
}

func HandleAction(c interfaces.GameCommunication, roomId, playerId uuid.UUID, data json.RawMessage) error {
	h, err := getGameType(roomId)
	if err != nil {
		return err
	}

	inGame, err := db.RoomHasUser(roomId, playerId)
	if err != nil {
		return err
	}

	if !inGame {
		return fmt.Errorf("you are not part of this game")
	}

	return h.HandleAction(c, roomId, playerId, data)
}

func HandleReady(c interfaces.GameCommunication, roomId, playerId uuid.UUID, data json.RawMessage) error {
	h, err := getGameType(roomId)
	if err != nil {
		return err
	}
	inGame, err := db.RoomHasUser(roomId, playerId)
	if err != nil {
		return err
	}

	if !inGame {
		return fmt.Errorf("you are not part of this game")
	}

	return h.HandleReady(c, roomId, playerId)
}

func New(roomId uuid.UUID) error {
	gameType, err := db.GetRoomProperty[string](roomId, "gameType")
	if err != nil {
		return err
	}

	h, ok := gameHandlers[gameType]
	if !ok {
		return fmt.Errorf("unrecognized game type")
	}

	playerMap, err := db.GetRoomUsers(roomId)
	if err != nil {
		return err
	}
	players := make([]uuid.UUID, len(playerMap))
	i := 0
	for k := range playerMap {
		players[i] = k
		i++
	}

	options, err := db.GetRoomProperty[json.RawMessage](roomId, "options")
	if err != nil {
		return fmt.Errorf("failed to get options from lobby")
	}

	err = h.New(roomId, options)
	if err != nil {
		return err
	}

	err = db.SetRoomProperty(roomId, "inGame", true)
	if err != nil {
		return err
	}

	return nil
}

func HandleLeave(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID) error {
	gameType, err := db.GetRoomProperty[string](gameId, "gameType")
	if err != nil {
		return fmt.Errorf("failed to get game type")
	}

	h, ok := gameHandlers[gameType]
	if !ok {
		return fmt.Errorf("unrecognized game type")
	}

	return h.HandleLeave(c, gameId, playerId)
}

func GetDefaultOptions(gameType string) (any, error) {
	h, ok := gameHandlers[gameType]
	if !ok {
		return nil, fmt.Errorf("unrecognized game type")
	}

	return h.DefaultOptions(), nil
}
