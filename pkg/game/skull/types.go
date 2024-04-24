package skull

import (
	"encoding/json"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/google/uuid"
)

type Options struct {
	DiscardRandom bool `json:"discardRandom"`
}

type GameState struct {
	Public  *PublicGameState  `json:"public"`
	Private *PrivateGameState `json:"private"`
}

type Action string

const a_place = "place"
const a_bid = "bid"
const a_pass = "pass"
const a_flip = "flip"

type ActionData struct {
	Option Action          `json:"option"`
	Data   json.RawMessage `json:"data"`
}

type ActionPlace struct {
	Tile Tile `json:"tile"`
}

type ActionBid struct {
	Bid int `json:"bid"`
}

type ActionFlip struct {
	Player uuid.UUID `json:"player"`
}

type Tile bool

const Rose Tile = false
const Skull Tile = true

var startingHand []Tile = []Tile{Rose, Rose, Rose, Skull}

type PublicGameState struct {
	TilesPlaced   map[uuid.UUID]int    `json:"tilesPlaced"`
	TilesRevealed map[uuid.UUID][]Tile `json:"tilesRevealed"`
	Bid           int                  `json:"bid"`
	Passed        []uuid.UUID          `json:"passed"`
	Points        map[uuid.UUID]int    `json:"points"`
	Flipper       uuid.UUID            `json:"flipper"`
	GameOver      bool                 `json:"gameOver"`
	TurnOrder     []uuid.UUID          `json:"turnOrder"`
	Turn          uuid.UUID            `json:"turn"`
}

type PrivateGameState struct {
	TilesPlaced []Tile `json:"tilesPlaced"`
	Tiles       []Tile `json:"tiles"`
}

const playerType = "player"

type DBProperty string

const d_bid DBProperty = "bid"
const d_gameOver DBProperty = "gameOver"
const d_flipper DBProperty = "flipper"
const d_currentTurn DBProperty = "currentTurn"
const d_passed DBProperty = "passed"

func GetProperty[T any](gameId uuid.UUID, p DBProperty) (T, error) {
	return db.GetGameProperty[T](gameId, string(p))
}

func SetProperty[T any](gameId uuid.UUID, p DBProperty, data T) error {
	return db.SetGameProperty(gameId, string(p), data)
}

type PlayerDBProperty string

const pd_tiles PlayerDBProperty = "tiles"
const pd_tilesPlaced PlayerDBProperty = "tilesPlaced"
const pd_tilesRevealed PlayerDBProperty = "tilesRevealed"
const pd_points PlayerDBProperty = "points"

func GetPlayerProperty[T any](gameId uuid.UUID, playerId uuid.UUID, p PlayerDBProperty) (T, error) {
	return db.GetPlayerProperty[T](gameId, playerId, string(p))
}

func SetPlayerProperty[T any](gameId uuid.UUID, playerId uuid.UUID, p PlayerDBProperty, data T) error {
	return db.SetPlayerProperty(gameId, playerId, string(p), data)
}
