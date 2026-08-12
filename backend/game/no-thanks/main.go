package nothanks

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/Bismyth/game-server/db"
	"github.com/Bismyth/game-server/interfaces"
	"github.com/google/uuid"
)

const Code = "nothanks"

const cacheExpireTime time.Duration = 2 * time.Hour

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) DefaultOptions() interface{} {
	return &Options{Rounds: 2}
}

func (h *Handler) New(gameId uuid.UUID, rawOptions []byte) error {
	var options Options
	err := json.Unmarshal(rawOptions, &options)
	if err != nil {
		return fmt.Errorf("failed to parse options")
	}

	players, err := db.GetRoomUserOrder(gameId)
	if err != nil {
		return err
	}

	if len(players) < 2 { //TODO: CHANGE BEFORE SHIPPING
		return fmt.Errorf("not enough players")
	}
	if len(players) > 7 {
		return fmt.Errorf("too many players")
	}

	// Randomise turn order
	rand.Shuffle(len(players), func(i, j int) {
		players[i], players[j] = players[j], players[i]
	})

	for _, p := range players {
		C_PLAYER.For(gameId).Add(p)
		PD_SCORE.MustSet(gameId, p, 0)
	}

	GD_GAME_OVER.MustSet(gameId, false)
	GD_ROUND.MustSet(gameId, 1)
	GD_TOTAL_ROUNDS.MustSet(gameId, options.Rounds)
	cachePrevious(gameId, nil)

	newRound(gameId)
	cachePublicGameState(gameId)

	return nil
}

func (h *Handler) HandleAction(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID, data json.RawMessage) error {
	var response ActionResponse

	var err error

	current := C_PLAYER.For(gameId).Current()
	if current != playerId {
		return fmt.Errorf("not your turn")
	}

	err = json.Unmarshal(data, &response)
	if err != nil {
		return fmt.Errorf("invalid player action")
	}

	switch response.Option {
	case ga_pass:
		err = handlePass(c, gameId, playerId)
	case ga_take:
		err = handleTake(c, gameId, playerId)
	default:
		err = fmt.Errorf("unrecognized player option")
	}
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) HandleReady(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID) error {
	publicGs := getPublicGameState(gameId)
	privateGs := getPrivateGameState(gameId, playerId)

	if publicGs.CurrentPlayer == playerId {
		c.ActionPrompt(playerId, allActions)
	}

	pr := getPrevious(gameId)
	if pr != nil {
		c.SendPlayer(playerId, pr)
	}
	c.SendPlayer(playerId, NewGameState(publicGs, privateGs))

	return nil
}

func (g *Handler) HandleLeave(c interfaces.GameCommunication, gameId, playerId uuid.UUID) error {
	// reset current round if remaining players > 3
	pTracker := C_PLAYER.For(gameId)
	isPlayer := pTracker.HasItem(playerId)
	if !isPlayer {
		return nil
	}
	err := pTracker.RemoveTarget(playerId)
	if err != nil {
		return err
	}
	if pTracker.Length() >= 3 {
		cachePrevious(gameId, nil)
		newRound(gameId)
		nextPlayer := C_PLAYER.For(gameId).Current()
		ps := cachePublicGameState(gameId)
		for _, player := range pTracker.GetAll() {
			p := loadPrivate(gameId, player)
			c.SendPlayer(player, NewGameState(ps, &p))
		}

		c.ActionPrompt(nextPlayer, allActions)
		return nil
	} else {
		return endGame(c, gameId, nil)
	}
}

func (h *Handler) Cleanup(gameId uuid.UUID) error {
	return cleanup(gameId)
}
