package skull

import (
	"encoding/json"
	"fmt"

	"github.com/Bismyth/game-server/db"
	"github.com/Bismyth/game-server/db/msg"
	"github.com/Bismyth/game-server/interfaces"
	"github.com/google/uuid"
)

func handleBid(c interfaces.GameCommunication, gameId, playerId uuid.UUID, data json.RawMessage) error {
	var bidData ActionBid
	err := json.Unmarshal(data, &bidData)
	if err != nil {
		return fmt.Errorf("invalid action data")
	}

	currentPassed := PGS_PASSED.MustGet(gameId)
	if isPassedPlayer(playerId, currentPassed) {
		return fmt.Errorf("you have already passed")
	}

	numPlayers := C_PLAYER.For(gameId).Length()
	if len(currentPassed) >= (int(numPlayers) - 1) {
		return fmt.Errorf("everyone else has passed, you must now flip")
	}
	totalTiles, err := totalTilesPlaced(gameId)
	if err != nil {
		return err
	}
	if totalTiles < bidData.Bid {
		return fmt.Errorf("cannot bid more than there are tiles")
	}

	currentBid := PGS_BID.MustGet(gameId)
	if currentBid >= bidData.Bid {
		return fmt.Errorf("must increase bid")
	}
	PGS_BID.MustSet(gameId, bidData.Bid)
	db.GameEvent(gameId, c).Log(msg.Msg().Player(playerId).Text(" has bid ").Bold(msg.Int(bidData.Bid)).String())

	if totalTiles == bidData.Bid {
		err = startFlipper(c, gameId, playerId)
		if err != nil {
			return err
		}
		return nil
	}

	_, err = goToUnpassedPlayer(gameId, currentPassed)
	if err != nil {
		return err
	}

	err = updatePublicGameState(c, gameId)
	if err != nil {
		return err
	}

	return nil
}

func handlePass(c interfaces.GameCommunication, gameId, playerId uuid.UUID, _ json.RawMessage) error {
	bid := PGS_BID.MustGet(gameId)
	if bid <= 0 {
		return fmt.Errorf("cannot pass without a bid")
	}

	currentPassed := PGS_PASSED.MustGet(gameId)
	if isPassedPlayer(playerId, currentPassed) {
		return fmt.Errorf("already passed")
	}

	currentPassed = append(currentPassed, playerId)
	PGS_PASSED.MustSet(gameId, currentPassed)

	db.GameEvent(gameId, c).Log(msg.Msg().Player(playerId).Text(" has passed").String())

	nextPlayer, err := goToUnpassedPlayer(gameId, currentPassed)
	if err != nil {
		return err
	}

	playerCount := C_PLAYER.For(gameId).Length()
	if len(currentPassed) >= (int(playerCount) - 1) {
		err = startFlipper(c, gameId, nextPlayer)
		if err != nil {
			return err
		}
		return nil
	}

	err = updatePublicGameState(c, gameId)
	if err != nil {
		return err
	}

	return nil
}
