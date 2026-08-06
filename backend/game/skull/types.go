package skull

import (
	"encoding/json"
	"errors"

	"github.com/Bismyth/game-server/db"
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
	TileCount     map[uuid.UUID]int    `json:"tileCount"`
	Bid           int                  `json:"bid"`
	Passed        []uuid.UUID          `json:"passed"`
	Points        map[uuid.UUID]int    `json:"points"`
	Flipper       uuid.UUID            `json:"flipper"`
	GameOver      bool                 `json:"gameOver"`
	TurnOrder     []uuid.UUID          `json:"turnOrder"`
	Turn          uuid.UUID            `json:"turn"`
	Round         int                  `json:"round"`
	PlayerLeft    bool                 `json:"playerLeft"`
}

func loadPublicGameState(gameId uuid.UUID) (PublicGameState, error) {
	var errs []error
	addErr := func(err error) {
		if err != nil {
			errs = append(errs, err)
		}
	}

	bid, err := PGS_BID.Get(gameId)
	addErr(err)
	gameOver, err := PGS_GAME_OVER.Get(gameId)
	addErr(err)
	passed, err := PGS_PASSED.Get(gameId)
	addErr(err)
	flipper, err := PGS_FLIPPER.Get(gameId)
	addErr(err)
	round, err := PGS_ROUND.Get(gameId)
	addErr(err)
	playerLeft, err := PGS_PLAYER_LEFT.Get(gameId)
	addErr(err)

	players, err := db.PlayerTypeGetAll(gameId, playerType)
	addErr(err)

	points, err := db.LoadPlayerProperty(gameId, players, "points", func(v int) (int, error) {
		return v, nil
	})
	addErr(err)

	tilesPlaced := make(map[uuid.UUID]int)
	tilesRevealed := make(map[uuid.UUID][]Tile)
	for _, p := range players {
		tp, err := PD_TILES_PLACED.Get(gameId, p)
		if err != nil {
			return PublicGameState{}, err
		}
		tilesPlaced[p] = len(tp)
		tr, err := PD_TILES_REVEALED.Get(gameId, p)
		if err != nil {
			return PublicGameState{}, err
		}
		tilesRevealed[p] = tp[max(len(tp)-tr, 0):]
	}
	addErr(err)
	tileCount, err := db.LoadAllProperty(gameId, players, PD_TILES, func(v []Tile) (int, error) {
		return len(v), nil
	})
	addErr(err)

	var turn uuid.UUID
	cursor := db.GetCursor(gameId, playerType)
	if cursor != nil {
		turn, err = cursor.Current()
		addErr(err)
	}

	if len(errs) > 0 {
		return PublicGameState{}, errors.Join(errs...)
	}
	return PublicGameState{
		Bid:           bid,
		GameOver:      gameOver,
		TurnOrder:     players,
		Points:        points,
		TilesPlaced:   tilesPlaced,
		TilesRevealed: tilesRevealed,
		TileCount:     tileCount,
		Passed:        passed,
		Flipper:       flipper,
		Round:         round,
		PlayerLeft:    playerLeft,
		Turn:          turn,
	}, nil
}

const playerType = "player"

var (
	PGS_BID         = db.Property[int]{Key: "bid"}
	PGS_GAME_OVER   = db.Property[bool]{Key: "gameOver"}
	PGS_FLIPPER     = db.Property[uuid.UUID]{Key: "flipper"}
	PGS_ROUND       = db.Property[int]{Key: "round"}
	PGS_PLAYER_LEFT = db.Property[bool]{Key: "playerLeft"}
	PGS_PASSED      = db.Property[[]uuid.UUID]{Key: "passed"}
)

var (
	PD_TILES          = db.PlayerProperty[[]Tile]{Key: "tiles"}
	PD_TILES_PLACED   = db.PlayerProperty[[]Tile]{Key: "tilesPlaced"}
	PD_TILES_REVEALED = db.PlayerProperty[int]{Key: "tilesRevealed"}
	PD_POINTS         = db.PlayerProperty[int]{Key: "points"}
)

type PrivateGameState struct {
	TilesPlaced []Tile `json:"tilesPlaced"`
	Tiles       []Tile `json:"tiles"`
}

func loadPrivateState(gameId uuid.UUID, playerId uuid.UUID) (PrivateGameState, error) {
	var errs []error
	addErr := func(err error) {
		if err != nil {
			errs = append(errs, err)
		}
	}

	tiles, err := PD_TILES.Get(gameId, playerId)
	addErr(err)
	tilesPlaced, err := PD_TILES_PLACED.Get(gameId, playerId)
	addErr(err)

	if len(errs) > 0 {
		return PrivateGameState{}, errors.Join(errs...)
	}

	return PrivateGameState{
		Tiles:       tiles,
		TilesPlaced: tilesPlaced,
	}, nil
}
