package liarsdice

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/Bismyth/game-server/db"
	"github.com/Bismyth/game-server/interfaces"
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

	players, err := db.GetRoomUserOrder(gameId)
	if err != nil {
		return err
	}

	if options.StartingDice <= 0 {
		return fmt.Errorf("must start game with more than 0 dice")
	}

	if options.StartingDice > 99 {
		return fmt.Errorf("too many dice")
	}

	if len(players) < 2 {
		return fmt.Errorf("not enough players")
	}

	if err := GD_BID.Set(gameId, ""); err != nil {
		return err
	}

	if err := GD_GAME_OVER.Set(gameId, false); err != nil {
		return err
	}

	// Randomise turn order
	rand.Shuffle(len(players), func(i, j int) {
		players[i], players[j] = players[j], players[i]
	})

	for _, player := range players {
		err := PD_DICE.Set(gameId, player, options.StartingDice)
		if err != nil {
			return err
		}
		db.PlayerGiveType(gameId, player, playerType)
	}

	pr := RoundInfo{
		Round: 0,
	}
	err = GD_PREVIOUS_ROUND.Set(gameId, pr)
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
	isPlayer := db.PlayerIsType(gameId, playerId, playerType)
	if !isPlayer {
		return nil
	}

	err := db.RemoveFromCursor(gameId, playerId, playerType)
	if err != nil {
		return err
	}

	end, err := checkEnd(gameId)
	if err != nil {
		return err
	}

	playerName, err := db.GetRoomUserName(gameId, playerId)
	if err != nil {
		return err
	}

	pr, err := generatePreviousRound(gameId, &ParsedRoundInfo{
		Leave: playerName,
	})
	if err != nil {
		return err
	}

	if end {
		err = endGame(c, gameId, pr)
		if err != nil {
			return err
		}
	} else {
		err = newRound(c, gameId, pr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Handler) Cleanup(gameId uuid.UUID) error {
	return cleanup(gameId)
}

func checkEnd(gameId uuid.UUID) (bool, error) {
	numPlayers, err := db.PlayerTypeCount(gameId, playerType)
	if err != nil {
		return false, err
	}

	return numPlayers <= 1, nil
}
