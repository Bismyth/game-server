package liarsdice

import (
	"fmt"

	"github.com/Bismyth/game-server/db"
	"github.com/Bismyth/game-server/db/msg"
	"github.com/Bismyth/game-server/interfaces"
	"github.com/google/uuid"
)

func getAllDice(gameId uuid.UUID) ([]int, error) {
	players, err := C_PLAYER.For(gameId).GetAll()
	if err != nil {
		return nil, err
	}
	hands, err := PD_HAND.GetMulti(gameId, players)
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
	currentBid := GD_BID.MustGet(gameId)

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

func currentPlayerLoseDice(c interfaces.GameCommunication, gameId uuid.UUID) (int, error) {
	playerId, err := C_PLAYER.For(gameId).Current()
	if err != nil {
		return 0, err
	}

	db.GameEvent(gameId, c).Log(msg.Msg().Player(playerId).Text(" has lost a dice").String())

	amount, err := PD_DICE.Get(gameId, playerId)
	if err != nil {
		return 0, err
	}

	newAmount := amount - 1

	err = PD_DICE.Set(gameId, playerId, newAmount)
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

	playerTracker := C_PLAYER.For(gameId)

	cu, err := playerTracker.Current()
	if err != nil {
		return err
	}
	pvInfo.CallUser = cu

	pv, err := playerTracker.PeekPrevious()
	if err != nil {
		return err
	}
	pvInfo.LastBid = pv

	if !bidRight {
		playerTracker.Previous()
	}
	db.GameEvent(gameId, c).Log(msg.Msg().Player(cu).Text(" has called out ").Player(pv).String())

	lostUser, err := playerTracker.Current()
	if err != nil {
		return err
	}
	pvInfo.DiceLost = lostUser

	pr, err := generatePreviousRound(gameId, &pvInfo)
	if err != nil {
		return err
	}

	a, err := currentPlayerLoseDice(c, gameId)
	if err != nil {
		return err
	}

	if a <= 0 {
		err := playerTracker.Remove()
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
		_, err := playerTracker.Next()
		if err != nil {
			return err
		}
	}

	if bidRight && a == 0 {
		_, err = playerTracker.Previous()
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
