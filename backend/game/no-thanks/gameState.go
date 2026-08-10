package nothanks

import (
	"github.com/Bismyth/game-server/db"
	"github.com/Bismyth/game-server/interfaces"
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

func updatePublicState(c interfaces.GameCommunication, gameId uuid.UUID) {
	pgs := cachePublicGameState(gameId)

	c.SendGlobal(GameState{
		Public: pgs,
	})
}

func getPrivateGameState(gameId uuid.UUID, playerId uuid.UUID) *PrivateGameState {
	if !C_PLAYER.For(gameId).HasItem(playerId) {
		return nil
	}

	ps := loadPrivate(gameId, playerId)

	return &ps
}
