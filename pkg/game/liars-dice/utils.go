package liarsdice

import (
	"math/rand/v2"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/Bismyth/game-server/pkg/interfaces"
	"github.com/google/uuid"
)

func rollHands(gameId uuid.UUID, players []uuid.UUID) error {
	for _, playerId := range players {
		numDice, err := db.GetPlayerProperty[int](gameId, playerId, "dice")
		if err != nil {
			return err
		}

		hand := make([]int, numDice)
		for i := range hand {
			hand[i] = rand.IntN(6) + 1
		}

		err = db.SetPlayerProperty(gameId, playerId, "hand", hand)
		if err != nil {
			return err
		}
	}

	return nil
}

func endGame(c interfaces.GameCommunication, gameId uuid.UUID) error {
	err := db.SetGameProperty(gameId, "gameOver", true)
	if err != nil {
		return err
	}

	err = cachePublicGameState(gameId)
	if err != nil {
		return err
	}

	err = db.ExpireCache(gameId, cacheExpireTime)
	if err != nil {
		return err
	}

	c.EndGame()

	pGs, err := getPublicGameState(gameId)
	if err != nil {
		return err
	}

	c.SendGlobal(GameState{
		Public: pGs,
	})

	err = cleanup(gameId)
	if err != nil {
		return err
	}

	return nil
}

func cleanup(gameId uuid.UUID) error {

	c := db.GetCursor(gameId, playerType)
	err := c.Delete()
	if err != nil {
		return err
	}

	err = db.DeletePlayerTypeList(gameId, playerType)
	if err != nil {
		return err
	}

	err = db.DeleteGame(gameId)
	if err != nil {
		return err
	}

	return nil
}
