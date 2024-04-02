package liarsdice

import (
	"fmt"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/Bismyth/game-server/pkg/interfaces"
	"github.com/google/uuid"
)

func getAllDice(gameId uuid.UUID) ([]int, error) {
	hands, err := db.GetMultiPlayerProperty[[]int](gameId, "hand")
	if err != nil {
		return nil, err
	}

	allDice := []int{}

	for _, h := range hands {
		allDice = append(allDice, h...)
	}
	return allDice, nil
}

func callRight(gameId uuid.UUID) (bool, error) {
	currentBid, err := db.GetGameProperty[string](gameId, "bid")
	if err != nil {
		return false, err
	}

	a, f, err := parseBid(currentBid)
	if err != nil {
		return false, err
	}

	allDice, err := getAllDice(gameId)
	if err != nil {
		return false, err
	}
	trueAmount := 0
	for _, dice := range allDice {
		if dice == 1 || dice == f {
			trueAmount++
		}
	}

	return trueAmount >= a, nil
}

func loseDice(gameId uuid.UUID, playerId uuid.UUID) error {
	amount, err := db.GetPlayerProperty[int](gameId, playerId, "dice")
	if err != nil {
		return err
	}

	newAmount := amount - 1

	err = db.SetPlayerProperty(gameId, playerId, "dice", newAmount)
	if err != nil {
		return err
	}

	if newAmount <= 0 {
		// TODO: Make player out

		// TODO: if only one player standing trigger game win
	}

	return nil
}

func getPreviousPlayer(gameId uuid.UUID) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func handleCall(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID) error {
	var err error

	cr, err := callRight(gameId)
	if err != nil {
		return fmt.Errorf("could not determine call")
	}

	if cr {
		err = loseDice(gameId, playerId)
	} else {
		oldPlayer, err := getPreviousPlayer(gameId)
		if err != nil {
			return err
		}
		err = loseDice(gameId, oldPlayer)
	}
	if err != nil {
		return err
	}
	players, err := db.GetGamePlayers(gameId)
	if err != nil {
		return err
	}
	err = rollHands(gameId, players)
	if err != nil {
		return err
	}

	for _, player := range players {
		privateGs, err := getPrivateGameState(gameId, player)
		if err != nil {
			return err
		}
		c.SendPlayer(player, GameState{Private: privateGs})
	}

	err = db.SetGameProperty(gameId, "bid", "")
	if err != nil {
		return err
	}

	return nil
}
