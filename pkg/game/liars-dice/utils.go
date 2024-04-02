package liarsdice

import (
	"math/rand/v2"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/google/uuid"
)

func incrementPlayerTurn(gameId uuid.UUID) error {
	numPlayers, err := db.GetGameProperty[int](gameId, "#players")
	if err != nil {
		return err
	}
	index, err := db.GetGameProperty[int](gameId, "turn")
	if err != nil {
		return err
	}

	index++
	index = index % numPlayers
	err = db.SetGameProperty(gameId, "turn", index)
	if err != nil {
		return err
	}

	players, err := db.GetGamePlayers(gameId)
	if err != nil {
		return err
	}

	err = db.SetGameProperty(gameId, "turnId", players[index])
	if err != nil {
		return err
	}

	return nil
}

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
