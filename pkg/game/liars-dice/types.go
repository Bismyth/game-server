package liarsdice

import "github.com/google/uuid"

type GameActions string

const ga_bid GameActions = "bid"
const ga_call GameActions = "call"

var allActions []GameActions = []GameActions{ga_bid, ga_call}

type GameState struct {
	Public  *PublicGameState  `json:"public"`
	Private *PrivateGameState `json:"private"`
}

type PublicGameState struct {
	HighestBid  string            `json:"highestBid"`
	PlayerTurn  uuid.UUID         `json:"playerTurn"`
	DiceAmounts map[uuid.UUID]int `json:"diceAmounts"`
	TurnOrder   []uuid.UUID       `json:"turnOrder"`
	GameOver    bool              `json:"gameOver"`
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
