package liarsdice

import (
	"fmt"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/Bismyth/game-server/pkg/interfaces"
	"github.com/google/uuid"
)

func getAllDice(gameId uuid.UUID) ([]int, error) {
	hands, err := db.GetMultiPlayerProperty[[]int](gameId, "hand", playerType)
	if err != nil {
		return nil, err
	}

	allDice := []int{}

	for _, h := range hands {
		allDice = append(allDice, h...)
	}
	return allDice, nil
}

// returns true if bid was met, false if bid was a lie
func evalBid(gameId uuid.UUID) (bool, error) {
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

func loseDiceAtCursor(gameId uuid.UUID, playerCursor *db.Cursor) (int, error) {

	playerId, err := playerCursor.Current()
	if err != nil {
		return 0, err
	}

	amount, err := db.GetPlayerProperty[int](gameId, playerId, "dice")
	if err != nil {
		return 0, err
	}

	newAmount := amount - 1

	err = db.SetPlayerProperty(gameId, playerId, "dice", newAmount)
	if err != nil {
		return 0, err
	}

	return newAmount, nil
}

func handleCall(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID) error {
	var err error

	bidRight, err := evalBid(gameId)
	if err != nil {
		return fmt.Errorf("could not determine call")
	}

	playerCursor := db.GetCursor(gameId, playerType)

	if !bidRight {
		playerCursor.Previous()
	}
	a, err := loseDiceAtCursor(gameId, playerCursor)
	if err != nil {
		return err
	}

	if a <= 0 {
		playerCursor.Remove()
		// TODO: if only one player standing trigger game win
	}

	if !bidRight && a > 0 {
		playerCursor.Next()
	}

	players, err := db.PlayerTypeGetAll(gameId, playerType)
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
