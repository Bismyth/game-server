package skull

import (
	"github.com/Bismyth/game-server/db"
	"github.com/google/uuid"
)

func cachePublicGameState(gameId uuid.UUID) (*PublicGameState, error) {
	pgs, err := loadPublicGameState(gameId)
	if err != nil {
		return nil, err
	}

	err = db.SetGameCache(gameId, pgs)
	if err != nil {
		return nil, err
	}

	return &pgs, nil
}

func getPublicGameState(gameId uuid.UUID) (*PublicGameState, error) {
	gs, err := db.GetGameCache[PublicGameState](gameId)
	if err != nil {
		return nil, err
	}

	return &gs, nil
}

func getPrivateGameState(gameId uuid.UUID, playerId uuid.UUID) (*PrivateGameState, error) {
	if !C_PLAYER.For(gameId).HasItem(playerId) {
		return nil, nil
	}

	ps, err := loadPrivateState(gameId, playerId)
	if err != nil {
		return nil, err
	}

	return &ps, nil
}
