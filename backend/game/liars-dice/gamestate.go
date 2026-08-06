package liarsdice

import (
	"github.com/Bismyth/game-server/db"
	"github.com/google/uuid"
)

func cachePublicGameState(gameId uuid.UUID) error {
	gs, err := loadPublicGameState(gameId)
	if err != nil {
		return err
	}

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
	if !db.PlayerIsType(gameId, playerId, playerType) {
		return nil, nil
	}

	hand, err := PD_HAND.Get(gameId, playerId)
	if err != nil {
		return nil, err
	}

	privateGs := PrivateGameState{Dice: hand}
	return &privateGs, nil
}
