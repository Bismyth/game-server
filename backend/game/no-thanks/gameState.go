package nothanks

import (
	"github.com/Bismyth/game-server/db"
	"github.com/google/uuid"
)

func cachePublicGameState(gameId uuid.UUID) *PublicGameState {
	pgs := loadPublic(gameId)

	err := db.SetGameCache(gameId, pgs)
	if err != nil {
		panic(err)
	}

	return &pgs
}

func getPublicGameState(gameId uuid.UUID) *PublicGameState {
	gs, err := db.GetGameCache[PublicGameState](gameId)
	if err != nil {
		panic(err)
	}

	return &gs
}

func cachePrevious(gameId uuid.UUID, previous *PreviousRound) {
	err := db.SetGameKey(gameId, "pr", previous)
	if err != nil {
		panic(err)
	}
}

func getPrevious(gameId uuid.UUID) *PreviousRound {
	pr, err := db.GetGameKey[PreviousRound](gameId, "pr")
	if err != nil {
		panic(err)
	}

	if pr.Type == "" {
		return nil
	}

	return &pr
}

func getPrivateGameState(gameId uuid.UUID, playerId uuid.UUID) *PrivateGameState {
	if !C_PLAYER.For(gameId).HasItem(playerId) {
		return nil
	}

	ps := loadPrivate(gameId, playerId)

	return &ps
}
