package liarsdice

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Bismyth/game-server/pkg/interfaces"
	"github.com/google/uuid"
)

func parseBid(bid string) (int, int, error) {
	returnErr := fmt.Errorf("could not parse bid")

	arr := strings.Split(bid, ",")
	if len(arr) != 2 {
		return 0, 0, returnErr
	}
	a, err := strconv.Atoi(arr[0])
	if err != nil {
		return 0, 0, returnErr
	}
	f, err := strconv.Atoi(arr[1])
	if err != nil {
		return 0, 0, returnErr
	}

	return a, f, nil
}

func checkValidBid(oldBid string, newBid string) bool {
	na, nf, err := parseBid(newBid)
	if err != nil {
		return false
	}

	if nf <= 1 || nf > 6 {
		return false
	}

	if na < 1 {
		return false
	}

	if oldBid == "" {
		return true
	}

	oa, of, err := parseBid(oldBid)
	if err != nil {
		return false
	}

	if na < oa {
		return false
	}

	if na == oa && nf <= of {
		return false
	}

	return true
}

func handleBid(c interfaces.GameCommunication, gameId uuid.UUID, bid string) error {
	oldBid, err := GetProperty[string](gameId, d_bid)
	if err != nil {
		return err
	}

	if !checkValidBid(oldBid, bid) {
		return fmt.Errorf("invalid bid")
	}

	err = SetProperty(gameId, d_bid, bid)
	if err != nil {
		return err
	}

	err = progressTurn(c, gameId)
	if err != nil {
		return err
	}

	return nil
}
