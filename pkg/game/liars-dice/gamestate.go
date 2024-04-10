package liarsdice

import (
	"github.com/Bismyth/game-server/pkg/db"
	"github.com/google/uuid"
)

func cachePublicGameState(gameId uuid.UUID) error {
	gs := PublicGameState{}

	players, err := db.PlayerTypeGetAll(gameId, playerType)
	if err != nil {
		return err
	}

	gs.TurnOrder = players

	hb, err := GetProperty[string](gameId, d_bid)
	if err != nil {
		return err
	}
	gs.HighestBid = hb

	gameOver, err := GetProperty[bool](gameId, d_gameOver)
	if err != nil {
		return err
	}
	gs.GameOver = gameOver

	if !gameOver {
		cursor := db.GetCursor(gameId, playerType)
		currentPLayer, err := cursor.Current()
		if err != nil {
			return err
		}
		gs.PlayerTurn = currentPLayer
	}

	playerDice := map[uuid.UUID]int{}

	for _, player := range players {
		num, err := db.GetPlayerProperty[int](gameId, player, "dice")
		if err != nil {
			return err
		}
		playerDice[player] = num
	}
	gs.DiceAmounts = playerDice

	pr, err := GetProperty[RoundInfo](gameId, d_previousRound)
	if err != nil {
		return err
	}
	gs.PreviousRound = pr

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

	hand, err := db.GetPlayerProperty[[]int](gameId, playerId, "hand")
	if err != nil {
		return nil, err
	}

	privateGs := PrivateGameState{Dice: hand}
	return &privateGs, nil
}
