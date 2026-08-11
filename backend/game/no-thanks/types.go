package nothanks

import (
	"github.com/Bismyth/game-server/db"
	"github.com/google/uuid"
)

type GameActions string

type Options struct {
	Rounds int `json:"rounds"`
}

const ga_pass GameActions = "pass"
const ga_take GameActions = "take"

var allActions []GameActions = []GameActions{ga_take, ga_pass}

type ActionResponse struct {
	Option GameActions
}

type GameState struct {
	Type    string            `json:"type"`
	Public  *PublicGameState  `json:"public"`
	Private *PrivateGameState `json:"private"`
}

func NewGameState(public *PublicGameState, private *PrivateGameState) GameState {
	return GameState{
		Type:    "state",
		Public:  public,
		Private: private,
	}
}

var DECK = db.Deck[int]{Key: "deck"}

var C_PLAYER = db.Cursor{Key: "player"}

var (
	GD_INPLAY         = db.Property[int]{Key: "inPlay"}
	GD_REMOVED        = db.Property[[]int]{Key: "removed"}
	GD_TOKENS_ON_CARD = db.Property[int]{Key: "tokensOnCard"}
	GD_TOTAL_ROUNDS   = db.Property[int]{Key: "totalRounds"}
	GD_ROUND          = db.Property[int]{Key: "round"}
	GD_GAME_OVER      = db.Property[bool]{Key: "gameOver"}
)

var (
	PD_CARDS  = db.PlayerProperty[[]int]{Key: "cards"}
	PD_TOKENS = db.PlayerProperty[int]{Key: "tokens"}
	PD_SCORE  = db.PlayerProperty[int]{Key: "score"}
)

type PublicGameState struct {
	InPlayCard    int                 `json:"inPlayCard"`
	TokensOnCard  int                 `json:"tokensOnCard"`
	PlayerCards   map[uuid.UUID][]int `json:"playerCards"`
	Score         map[uuid.UUID]int   `json:"score"`
	DeckLeft      int                 `json:"deckLeft"`
	TurnOrder     []uuid.UUID         `json:"turnOrder"`
	CurrentPlayer uuid.UUID           `json:"currentPlayer"`
	CurrentRound  int                 `json:"currentRound"`
	TotalRounds   int                 `json:"totalRounds"`
}

type PreviousRound struct {
	Type        string              `json:"type"`
	Score       map[uuid.UUID]int   `json:"score"`
	PlayerCards map[uuid.UUID][]int `json:"playerCards"`
	Removed     []int               `json:"removed"`
	Round       int                 `json:"round"`
}

func loadPublic(gameId uuid.UUID) PublicGameState {

	inPlay := GD_INPLAY.MustGet(gameId)
	tOnCard := GD_TOKENS_ON_CARD.MustGet(gameId)
	round := GD_ROUND.MustGet(gameId)
	totalRounds := GD_TOTAL_ROUNDS.MustGet(gameId)

	deckSize := int(DECK.For(gameId).Length())

	players := C_PLAYER.For(gameId).GetAll()

	playerCards, err := PD_CARDS.GetMap(gameId, players)
	if err != nil {
		panic(err)
	}
	score, err := PD_SCORE.GetMap(gameId, players)
	if err != nil {
		panic(err)
	}

	current := C_PLAYER.For(gameId).Current()

	return PublicGameState{
		InPlayCard:    inPlay,
		TokensOnCard:  tOnCard,
		DeckLeft:      deckSize,
		PlayerCards:   playerCards,
		Score:         score,
		TurnOrder:     players,
		CurrentPlayer: current,
		CurrentRound:  round,
		TotalRounds:   totalRounds,
	}
}

type PrivateGameState struct {
	Tokens int `json:"tokens"`
}

func loadPrivate(gameId, playerId uuid.UUID) PrivateGameState {
	tokens := PD_TOKENS.MustGet(gameId, playerId)
	return PrivateGameState{
		Tokens: tokens,
	}
}
