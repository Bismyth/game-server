package liarsdice

import (
	"encoding/json"
	"fmt"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/Bismyth/game-server/pkg/interfaces"
	"github.com/google/uuid"
)

const Code = "liarsdice"

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) New(gameId uuid.UUID, rawOptions json.RawMessage) error {
	var options Options
	err := json.Unmarshal(rawOptions, &options)
	if err != nil {
		return fmt.Errorf("failed to parse options")
	}

	if err := db.SetGameProperty(gameId, "bid", ""); err != nil {
		return err
	}

	players, err := db.GetGamePlayers(gameId)
	if err != nil {
		return err
	}

	turnIndex := 0

	if err := db.SetGameProperty(gameId, "turn", turnIndex); err != nil {
		return err
	}
	if err := db.SetGameProperty(gameId, "turnId", players[turnIndex]); err != nil {
		return err
	}

	for _, player := range players {
		if err := db.SetPlayerProperty(gameId, player, "dice", options.StartingDice); err != nil {
			return err
		}
	}

	rollHands(gameId, players)

	err = cachePublicGameState(gameId)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) HandleAction(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID, data json.RawMessage) error {

	var response ActionResponse

	var err error

	err = json.Unmarshal(data, &response)
	if err != nil {
		return fmt.Errorf("invalid player action")
	}

	switch response.Option {
	case ga_bid:
		err = handleBid(c, gameId, playerId, response.Data.Bid)
	case ga_call:
		err = handleCall(c, gameId, playerId)
	default:
		err = fmt.Errorf("unrecognized player option")
	}

	if err != nil {
		return err
	}

	err = incrementPlayerTurn(gameId)
	if err != nil {
		return err
	}
	err = cachePublicGameState(gameId)
	if err != nil {
		return err
	}

	publicGs, err := getPublicGameState(gameId)
	if err != nil {
		return err
	}
	gs := GameState{Public: publicGs}

	activePlayer, err := db.GetGameProperty[uuid.UUID](gameId, "turnId")
	if err != nil {
		return err
	}

	c.SendGlobal(gs)
	c.ActionPrompt(activePlayer, allActions)

	return nil
}

func (h *Handler) HandleReady(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID) error {
	publicGs, err := getPublicGameState(gameId)
	if err != nil {
		return err
	}
	privateGs, err := getPrivateGameState(gameId, playerId)
	if err != nil {
		return err
	}

	activePlayer, err := db.GetGameProperty[uuid.UUID](gameId, "turnId")
	if err != nil {
		return err
	}

	if activePlayer == playerId {
		c.ActionPrompt(playerId, allActions)
	}

	c.SendPlayer(playerId, GameState{
		Public:  publicGs,
		Private: privateGs,
	})

	return nil
}
