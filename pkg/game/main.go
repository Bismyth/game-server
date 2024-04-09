package game

import (
	"encoding/json"
	"fmt"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/Bismyth/game-server/pkg/interfaces"
	"github.com/google/uuid"

	liarsdice "github.com/Bismyth/game-server/pkg/game/liars-dice"
)

type GameHandler interface {
	HandleAction(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID, data json.RawMessage) error
	HandleReady(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID) error
	New(gameId uuid.UUID, options []byte) error
	DefaultOptions() interface{}
	HandleLeave(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID) error
}

var gameHandlers map[string]GameHandler = map[string]GameHandler{
	liarsdice.Code: liarsdice.New(),
}

type SharedState struct {
	Id uuid.UUID
}

func getGameType(data json.RawMessage) (GameHandler, uuid.UUID, error) {
	var state SharedState
	err := json.Unmarshal(data, &state)
	if err != nil {
		return nil, uuid.Nil, fmt.Errorf("failed to unmarshal gameId")
	}

	gameType, err := db.GetLobbyProperty[string](state.Id, "gameType")
	if err != nil {
		return nil, state.Id, err
	}

	h, ok := gameHandlers[gameType]
	if !ok {
		return nil, state.Id, fmt.Errorf("game type not found in list")
	}
	return h, state.Id, nil
}

func HandleAction(c interfaces.GameCommunication, playerId uuid.UUID, data json.RawMessage) error {
	h, id, err := getGameType(data)
	if err != nil {
		return err
	}

	inGame, err := db.IsUserInLobby(id, playerId)
	if err != nil {
		return err
	}

	if !inGame {
		return fmt.Errorf("you are not part of this game")
	}

	return h.HandleAction(c, id, playerId, data)
}

func HandleReady(c interfaces.GameCommunication, playerId uuid.UUID, data json.RawMessage) error {
	h, id, err := getGameType(data)
	if err != nil {
		return err
	}
	inGame, err := db.IsUserInLobby(id, playerId)
	if err != nil {
		return err
	}

	if !inGame {
		return fmt.Errorf("you are not part of this game")
	}

	return h.HandleReady(c, id, playerId)
}

func New(lobbyId uuid.UUID) error {
	gameType, err := db.GetLobbyProperty[string](lobbyId, "gameType")
	if err != nil {
		return err
	}

	h, ok := gameHandlers[gameType]
	if !ok {
		return fmt.Errorf("unrecognized game type")
	}

	playerMap, err := db.GetLobbyUsers(lobbyId)
	if err != nil {
		return err
	}
	players := make([]uuid.UUID, len(playerMap))
	i := 0
	for k := range playerMap {
		players[i] = k
		i++
	}

	options, err := db.GetLobbyProperty[json.RawMessage](lobbyId, "options")
	if err != nil {
		return fmt.Errorf("failed to get options from lobby")
	}

	err = h.New(lobbyId, options)
	if err != nil {
		return err
	}

	err = db.SetLobbyProperty(lobbyId, "inGame", true)
	if err != nil {
		return err
	}

	return nil
}

func HandleLeave(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID) error {
	gameType, err := db.GetLobbyProperty[string](gameId, "gameType")
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
