package liarsdice

import (
	"github.com/Bismyth/game-server/pkg/db"
	"github.com/google/uuid"
)

func cachePublicGameState(gameId uuid.UUID) error {
	gs := PublicGameState{}

	players, err := db.GetGamePlayers(gameId)
	if err != nil {
		return err
	}

	hb, err := db.GetGameProperty[string](gameId, "bid")
	if err != nil {
		return err
	}
	gs.HighestBid = hb
	turnIndex, err := db.GetGameProperty[int](gameId, "turn")
	if err != nil {
		return err
	}
	gs.PlayerTurn = players[turnIndex]

	playerDice := map[uuid.UUID]int{}

	for _, player := range players {
		num, err := db.GetPlayerProperty[int](gameId, player, "dice")
		if err != nil {
			return err
		}
		playerDice[player] = num
	}
	gs.DiceAmounts = playerDice

	err = db.SetGameProperty(gameId, "cache", gs)
	if err != nil {
		return err
	}

	return nil
}

func getPublicGameState(gameId uuid.UUID) (*PublicGameState, error) {
	gs, err := db.GetGameProperty[PublicGameState](gameId, "cache")
	if err != nil {
		return nil, err
	}

	return &gs, nil
}

func getPrivateGameState(gameId uuid.UUID, playerId uuid.UUID) (*PrivateGameState, error) {
	hand, err := db.GetPlayerProperty[[]int](gameId, playerId, "hand")
	if err != nil {
		return nil, err
	}

	privateGs := PrivateGameState{Dice: hand}
	return &privateGs, nil
}
