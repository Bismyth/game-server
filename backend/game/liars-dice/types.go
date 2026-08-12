package liarsdice

import (
	"errors"

	"github.com/Bismyth/game-server/db"
	"github.com/google/uuid"
)

type GameActions string

const ga_bid GameActions = "bid"
const ga_call GameActions = "call"

var allActions []GameActions = []GameActions{ga_bid, ga_call}

type GameState struct {
	Public  *PublicGameState  `json:"public"`
	Private *PrivateGameState `json:"private"`
}

type PublicGameState struct {
	HighestBid    string            `json:"highestBid"`
	PlayerTurn    uuid.UUID         `json:"playerTurn"`
	DiceAmounts   map[uuid.UUID]int `json:"diceAmounts"`
	TurnOrder     []uuid.UUID       `json:"turnOrder"`
	GameOver      bool              `json:"gameOver"`
	PreviousRound RoundInfo         `json:"previousRound"`
}

type RoundInfo struct {
	Round      int                 `json:"round"`
	HighestBid string              `json:"highestBid"`
	Hands      map[uuid.UUID][]int `json:"hands"`
	LastBid    uuid.UUID           `json:"lastBid"`
	CallUser   uuid.UUID           `json:"callUser"`
	DiceLost   uuid.UUID           `json:"diceLost"`
	Leave      string              `json:"leave"`
}

type ParsedRoundInfo struct {
	LastBid  uuid.UUID
	CallUser uuid.UUID
	DiceLost uuid.UUID
	Leave    string
}

type PrivateGameState struct {
	Dice []int `json:"dice"`
}

type ActionResponse struct {
	Option GameActions
	Data   struct {
		Bid string
	}
}

type Options struct {
	StartingDice int `json:"startingDice"`
}

var (
	GD_BID            = db.Property[string]{Key: "bid"}
	GD_PREVIOUS_ROUND = db.Property[RoundInfo]{Key: "previousRound"}
	GD_GAME_OVER      = db.Property[bool]{Key: "gameOver"}
)

func loadPublicGameState(gameId uuid.UUID) (PublicGameState, error) {
	var errs []error
	addErr := func(err error) {
		if err != nil {
			errs = append(errs, err)
		}
	}
	pTracker := C_PLAYER.For(gameId)
	players := pTracker.GetAll()
	hb, err := GD_BID.Get(gameId)
	addErr(err)
	gameOver, err := GD_GAME_OVER.Get(gameId)
	addErr(err)

	currentPlayer := uuid.Nil
	if !gameOver {
		currentPlayer = pTracker.Current()
	}

	playerDice, err := PD_DICE.GetMap(gameId, players)
	addErr(err)

	pr, err := GD_PREVIOUS_ROUND.Get(gameId)
	addErr(err)

	if len(errs) > 0 {
		return PublicGameState{}, errors.Join(errs...)
	}

	return PublicGameState{
		TurnOrder:     players,
		HighestBid:    hb,
		GameOver:      gameOver,
		PlayerTurn:    currentPlayer,
		DiceAmounts:   playerDice,
		PreviousRound: pr,
	}, nil
}

var (
	PD_DICE = db.PlayerProperty[int]{Key: "dice"}
	PD_HAND = db.PlayerProperty[[]int]{Key: "hand"}
)

var C_PLAYER = db.Cursor{Key: "player"}
