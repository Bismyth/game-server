package liarsdice

import (
	"fmt"

	"github.com/Bismyth/game-server/db"
	"github.com/Bismyth/game-server/db/msg"
	"github.com/Bismyth/game-server/interfaces"
	"github.com/google/uuid"
)

func getAllDice(gameId uuid.UUID) ([]int, error) {
	players := C_PLAYER.For(gameId).GetAll()

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

func currentPlayerLoseDice(c interfaces.GameCommunication, gameId uuid.UUID) int {
	playerId := C_PLAYER.For(gameId).Current()
	db.GameEvent(gameId, c).Log(msg.Msg().Player(playerId).Text(" has lost a dice").String())

	amount := PD_DICE.MustGet(gameId, playerId)
	newAmount := amount - 1

	PD_DICE.MustSet(gameId, playerId, newAmount)

	return newAmount
}

func handleCall(c interfaces.GameCommunication, gameId uuid.UUID) error {
	var err error
	var pvInfo ParsedRoundInfo

	bidRight, err := evalBid(gameId)
	if err != nil {
		return fmt.Errorf("could not determine call")
	}

	playerTracker := C_PLAYER.For(gameId)

	cu := playerTracker.Current()
	pvInfo.CallUser = cu

	pv := playerTracker.PeekPrevious()
	pvInfo.LastBid = pv

	if !bidRight {
		playerTracker.Previous()
	}
	db.GameEvent(gameId, c).Log(msg.Msg().Player(cu).Text(" has called out ").Player(pv).String())

	lostUser := playerTracker.Current()
	pvInfo.DiceLost = lostUser

	pr, err := generatePreviousRound(gameId, &pvInfo)
	if err != nil {
		return err
	}

	a := currentPlayerLoseDice(c, gameId)
	if a <= 0 {
		playerTracker.Remove()
		if checkEnd(gameId) {
			endGame(c, gameId, pr)
			return nil
		}
	}

	if !bidRight && a > 0 {
		playerTracker.Next()
	}

	if bidRight && a == 0 {
		playerTracker.Previous()
	}

	err = newRound(c, gameId, pr)
	if err != nil {
		return err
	}

	return nil
}
