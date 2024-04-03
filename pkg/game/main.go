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
	New(gameId uuid.UUID, options json.RawMessage) error
}

var gameHandlers map[string]GameHandler = map[string]GameHandler{
	liarsdice.Code: liarsdice.New(),
}

type SharedState struct {
	Id   uuid.UUID
	Type string
}

func getGameType(data json.RawMessage) (GameHandler, uuid.UUID, error) {
	var state SharedState
	err := json.Unmarshal(data, &state)
	if err != nil {
		return nil, uuid.Nil, fmt.Errorf("failed to determine game type")
	}

	h, ok := gameHandlers[state.Type]
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

	if !db.IsUserInGame(playerId, id) {
		return fmt.Errorf("you are not part of this game")
	}

	turnId, err := db.GetGameProperty[uuid.UUID](id, "turnId")
	if err != nil || turnId != playerId {
		return fmt.Errorf("it is not your turn")
	}

	return h.HandleAction(c, id, playerId, data)
}

func HandleReady(c interfaces.GameCommunication, playerId uuid.UUID, data json.RawMessage) error {
	h, id, err := getGameType(data)
	if err != nil {
		return err
	}
	if !db.IsUserInGame(playerId, id) {
		return fmt.Errorf("you are not part of this game")
	}

	return h.HandleReady(c, id, playerId)
}

type NewGame struct {
	Type    string
	LobbyId uuid.UUID
	Options json.RawMessage
}

func New(data json.RawMessage) (uuid.UUID, error) {
	var newGame NewGame
	err := json.Unmarshal(data, &newGame)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to unmarshal new game data")
	}

	h, ok := gameHandlers[newGame.Type]
	if !ok {
		return uuid.Nil, fmt.Errorf("unrecognized game type")
	}

	playerMap, err := db.GetLobbyUsers(newGame.LobbyId)
	if err != nil {
		return uuid.Nil, err
	}
	players := make([]uuid.UUID, len(playerMap))
	i := 0
	for k := range playerMap {
		players[i] = k
		i++
	}

	id, err := db.NewGame(newGame.LobbyId, newGame.Type, players)
	if err != nil {
		return id, err
	}

	err = h.New(id, newGame.Options)
	if err != nil {
		return id, err
	}

	return id, nil
}
