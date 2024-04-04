package liarsdice

import (
	"encoding/json"
	"fmt"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/Bismyth/game-server/pkg/interfaces"
	"github.com/google/uuid"
)

const Code = "liarsdice"

const playerType = "player"

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) New(gameId uuid.UUID, rawOptions []byte) error {
	var options Options
	err := json.Unmarshal(rawOptions, &options)
	if err != nil {
		return fmt.Errorf("failed to parse options")
	}

	if err := db.SetGameProperty(gameId, "bid", ""); err != nil {
		return err
	}

	players, err := db.GetLobbyUserIds(gameId)
	if err != nil {
		return err
	}

	for _, player := range players {
		if err := db.SetPlayerProperty(gameId, player, "dice", options.StartingDice); err != nil {
			return err
		}
		db.PlayerGiveType(gameId, player, playerType)
	}

	c := db.GetCursor(gameId, playerType)
	c.Reset()

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

	cursor := db.GetCursor(gameId, playerType)
	current, err := cursor.Current()
	if err != nil {
		return err
	}
	if current != playerId {
		return fmt.Errorf("not your turn")
	}

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

	nextPlayer, err := cursor.Next()
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

	c.SendGlobal(gs)
	c.ActionPrompt(nextPlayer, allActions)

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

	cursor := db.GetCursor(gameId, playerType)
	activePlayer, err := cursor.Current()
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
