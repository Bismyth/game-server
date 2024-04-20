package skull

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/Bismyth/game-server/pkg/interfaces"
	"github.com/google/uuid"
)

const Code = "skull"

const cacheExpireTime time.Duration = 2 * time.Hour

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) DefaultOptions() interface{} {
	return &Options{
		DiscardRandom: true,
	}
}

func (h *Handler) New(gameId uuid.UUID, rawOptions []byte) (err error) {
	defer func() {
		if err != nil {
			e := cleanup(gameId)
			if e != nil {
				log.Printf("failed to cleanup game: %s", gameId.String())
			}
		}
	}()

	var options Options

	err = json.Unmarshal(rawOptions, &options)
	if err != nil {
		return
	}

	err = SetProperty(gameId, d_gameOver, false)
	if err != nil {
		return
	}

	err = SetProperty(gameId, d_bid, 0)
	if err != nil {
		return
	}

	err = SetProperty(gameId, d_flipper, uuid.Nil)
	if err != nil {
		return
	}

	players, err := db.GetRoomUserOrder(gameId)
	if err != nil {
		return
	}

	if len(players) < 3 {
		err = fmt.Errorf("not enough players")
		return
	}

	if len(players) > 6 {
		err = fmt.Errorf("too many players")
		return
	}

	// Randomise turn order
	rand.Shuffle(len(players), func(i, j int) {
		players[i], players[j] = players[j], players[i]
	})

	for _, player := range players {
		err = db.PlayerGiveType(gameId, player, playerType)
		if err != nil {
			return
		}
		err = SetPlayerProperty(gameId, player, pd_points, 0)
		if err != nil {
			return
		}
		err = SetPlayerProperty(gameId, player, pd_tiles, startingHand)
		if err != nil {
			return
		}
		err = SetPlayerProperty(gameId, player, pd_tilesPlaced, []Tile{})
		if err != nil {
			return
		}
		err = SetPlayerProperty(gameId, player, pd_tilesRevealed, 0)
		if err != nil {
			return
		}
	}

	c := db.GetCursor(gameId, playerType)
	c.Reset()

	err = cachePublicGameState(gameId)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) HandleAction(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID, data json.RawMessage) error {
	return fmt.Errorf("not implemented")
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

	if len(privateGs.TilesPlaced) <= 0 {
		c.ActionPrompt(playerId, []Action{a_initialPlace})
	}

	c.SendPlayer(playerId, GameState{
		Public:  publicGs,
		Private: privateGs,
	})

	return nil
}

func (h *Handler) HandleLeave(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID) error {
	return fmt.Errorf("not implemented")
}

func (h *Handler) Cleanup(gameId uuid.UUID) error {
	return cleanup(gameId)
}
