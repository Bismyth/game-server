package liarsdice

import (
	"github.com/Bismyth/game-server/pkg/db"
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
}

type ParsedRoundInfo struct {
	LastBid  uuid.UUID
	CallUser uuid.UUID
	DiceLost uuid.UUID
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

type DBProperty string

const d_bid DBProperty = "bid"
const d_previousRound DBProperty = "previousRound"
const d_gameOver DBProperty = "gameOver"

func GetProperty[T any](gameId uuid.UUID, p DBProperty) (T, error) {
	return db.GetGameProperty[T](gameId, string(p))
}

func SetProperty[T any](gameId uuid.UUID, p DBProperty, data T) error {
	return db.SetGameProperty(gameId, string(p), data)
}
