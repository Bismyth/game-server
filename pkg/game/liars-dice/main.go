package liarsdice

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/Bismyth/game-server/pkg/interfaces"
	"github.com/google/uuid"
)

const Code = "liarsdice"

const playerType = "player"

const cacheExpireTime time.Duration = 2 * time.Hour

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) DefaultOptions() interface{} {
	return &Options{
		StartingDice: 5,
	}
}

func (h *Handler) New(gameId uuid.UUID, rawOptions []byte) error {
	var options Options
	err := json.Unmarshal(rawOptions, &options)
	if err != nil {
		return fmt.Errorf("failed to parse options")
	}

	if err := SetProperty(gameId, d_bid, ""); err != nil {
		return err
	}

	if err := SetProperty(gameId, d_gameOver, false); err != nil {
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

	pr := RoundInfo{
		Round: 0,
	}
	err = SetProperty(gameId, d_previousRound, pr)
	if err != nil {
		return err
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
		err = handleBid(c, gameId, response.Data.Bid)
	case ga_call:
		err = handleCall(c, gameId)
	default:
		err = fmt.Errorf("unrecognized player option")
	}
	if err != nil {
		return err
	}

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

	if publicGs.PlayerTurn == playerId {
		c.ActionPrompt(playerId, allActions)
	}

	c.SendPlayer(playerId, GameState{
		Public:  publicGs,
		Private: privateGs,
	})

	return nil
}

func (h *Handler) HandleLeave(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID) error {
	cursor := db.GetCursor(gameId, playerType)
	current, err := cursor.Current()
	if err != nil {
		return err
	}
	if current == playerId {
		err := cursor.Remove()
		if err != nil {
			return err
		}
	} else {
		err := cursor.SeekIndex(playerId)
		if err != nil {
			return err
		}
		err = cursor.Remove()
		if err != nil {
			return err
		}
		err = cursor.SeekIndex(current)
		if err != nil {
			return err
		}
	}

	end, err := checkEnd(gameId)
	if err != nil {
		return err
	}

	if end {
		err = endGame(c, gameId, nil)
		if err != nil {
			return err
		}
	} else {
		err = cachePublicGameState(gameId)
		if err != nil {
			return err
		}

		pGs, err := getPublicGameState(gameId)
		if err != nil {
			return err
		}

		c.SendGlobal(GameState{
			Public: pGs,
		})
	}

	return nil
}

func checkEnd(gameId uuid.UUID) (bool, error) {
	numPlayers, err := db.PlayerTypeCount(gameId, playerType)
	if err != nil {
		return false, err
	}

	return numPlayers <= 1, nil
}
