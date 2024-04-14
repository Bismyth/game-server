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
	currentBid, err := GetProperty[string](gameId, d_bid)
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

func handleCall(c interfaces.GameCommunication, gameId uuid.UUID) error {
	var err error
	var pvInfo ParsedRoundInfo

	bidRight, err := evalBid(gameId)
	if err != nil {
		return fmt.Errorf("could not determine call")
	}

	playerCursor := db.GetCursor(gameId, playerType)

	cu, err := playerCursor.Current()
	if err != nil {
		return err
	}
	pvInfo.CallUser = cu

	pv, err := playerCursor.PeekPrevious()
	if err != nil {
		return err
	}
	pvInfo.LastBid = pv

	if !bidRight {
		playerCursor.Previous()
	}

	lostUser, err := playerCursor.Current()
	if err != nil {
		return err
	}
	pvInfo.DiceLost = lostUser

	pr, err := generatePreviousRound(gameId, &pvInfo)
	if err != nil {
		return err
	}

	a, err := loseDiceAtCursor(gameId, playerCursor)
	if err != nil {
		return err
	}

	if a <= 0 {
		err := playerCursor.Remove()
		if err != nil {
			return err
		}
		end, err := checkEnd(gameId)
		if err != nil {
			return err
		}
		if end {
			endGame(c, gameId, pr)
			return nil
		}
	}

	if !bidRight && a > 0 {
		_, err := playerCursor.Next()
		if err != nil {
			return err
		}
	}

	if bidRight && a == 0 {
		_, err = playerCursor.Previous()
		if err != nil {
			return err
		}
	}

	err = newRound(c, gameId, pr)
	if err != nil {
		return err
	}

	return nil
}
