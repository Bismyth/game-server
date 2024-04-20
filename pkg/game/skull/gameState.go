package skull

import (
	"github.com/Bismyth/game-server/pkg/db"
	"github.com/google/uuid"
)

func cachePublicGameState(gameId uuid.UUID) error {
	gs := PublicGameState{}

	bid, err := GetProperty[int](gameId, d_bid)
	if err != nil {
		return err
	}
	gs.Bid = bid

	flipper, err := GetProperty[uuid.UUID](gameId, d_flipper)
	if err != nil {
		return err
	}
	gs.Flipper = flipper

	gameOver, err := GetProperty[bool](gameId, d_gameOver)
	if err != nil {
		return err
	}
	gs.GameOver = gameOver

	tilesPlaced := make(map[uuid.UUID]int)
	tilesRevealed := make(map[uuid.UUID][]Tile)
	points := make(map[uuid.UUID]int)

	players, err := db.PlayerTypeGetAll(gameId, playerType)
	if err != nil {
		return err
	}

	gs.TurnOrder = players

	for _, playerId := range players {
		tp, err := GetPlayerProperty[[]Tile](gameId, playerId, pd_tilesPlaced)
		if err != nil {
			return err
		}
		tilesPlaced[playerId] = len(tp)

		tr, err := GetPlayerProperty[int](gameId, playerId, pd_tilesRevealed)
		if err != nil {
			return err
		}
		tilesRevealed[playerId] = tp[max(len(tp)-tr, 0):]

		p, err := GetPlayerProperty[int](gameId, playerId, pd_tilesRevealed)
		if err != nil {
			return err
		}
		points[playerId] = p
	}
	gs.TilesPlaced = tilesPlaced
	gs.TilesRevealed = tilesRevealed
	gs.Points = points

	err = db.SetGameCache(gameId, gs)
	if err != nil {
		return err
	}

	return nil
}

func getPublicGameState(gameId uuid.UUID) (*PublicGameState, error) {
	gs, err := db.GetGameCache[PublicGameState](gameId)
	if err != nil {
		return nil, err
	}

	return &gs, nil
}

func getPrivateGameState(gameId uuid.UUID, playerId uuid.UUID) (*PrivateGameState, error) {
	ps := PrivateGameState{}

	tp, err := GetPlayerProperty[[]Tile](gameId, playerId, pd_tilesPlaced)
	if err != nil {
		return nil, err
	}
	ps.TilesPlaced = tp

	t, err := GetPlayerProperty[[]Tile](gameId, playerId, pd_tiles)
	if err != nil {
		return nil, err
	}
	ps.Tiles = t

	return &ps, nil
}
