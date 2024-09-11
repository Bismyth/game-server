package skull

import (
	"encoding/json"
	"fmt"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/Bismyth/game-server/pkg/interfaces"
	"github.com/google/uuid"
)

func handleBid(c interfaces.GameCommunication, gameId, playerId uuid.UUID, data json.RawMessage) error {

	turn, err := GetProperty[uuid.UUID](gameId, d_currentTurn)
	if err != nil {
		return err
	}
	if turn != playerId {
		return fmt.Errorf("not your turn")
	}

	var bidData ActionBid
	err = json.Unmarshal(data, &bidData)
	if err != nil {
		return fmt.Errorf("invalid action data")
	}

	tilesPlaced, err := countTilesPlaced(gameId)
	if err != nil {
		return err
	}
	if !allPlayersTilesPlaced(tilesPlaced) {
		return fmt.Errorf("all players must place at least one tile")
	}

	currentPassed, err := GetProperty[[]uuid.UUID](gameId, d_passed)
	if err != nil {
		return err
	}

	if isPassedPlayer(playerId, currentPassed) {
		return fmt.Errorf("you have already passed")
	}

	numPlayers, err := db.PlayerTypeCount(gameId, playerType)
	if err != nil {
		return err
	}

	if len(currentPassed) >= (int(numPlayers) - 1) {
		return fmt.Errorf("everyone else has passed, you must now flip")
	}

	if totalTilesPlaced(tilesPlaced) < bidData.Bid {
		return fmt.Errorf("cannot bid more than there are tiles")
	}

	currentBid, err := GetProperty[int](gameId, d_bid)
	if err != nil {
		return err
	}
	if currentBid <= bidData.Bid {
		return fmt.Errorf("must increase bid")
	}

	err = SetProperty(gameId, d_bid, bidData.Bid)
	if err != nil {
		return err
	}

	nextPlayer, err := findNextUnpassedPlayer(gameId, currentPassed)
	if err != nil {
		return err
	}

	err = SetProperty(gameId, d_currentTurn, nextPlayer)
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
	turn, err := GetProperty[uuid.UUID](gameId, d_currentTurn)
	if err != nil {
		return err
	}
	if turn != playerId {
		return fmt.Errorf("not your turn")
	}

	currentPassed, err := GetProperty[[]uuid.UUID](gameId, d_passed)
	if err != nil {
		return err
	}

	if isPassedPlayer(playerId, currentPassed) {
		return fmt.Errorf("already passed")
	}

	currentPassed = append(currentPassed, playerId)
	err = SetProperty(gameId, d_passed, currentPassed)
	if err != nil {
		return err
	}

	nextPlayer, err := findNextUnpassedPlayer(gameId, currentPassed)
	if err != nil {
		return err
	}
	err = SetProperty(gameId, d_currentTurn, nextPlayer)
	if err != nil {
		return err
	}

	playerCount, err := db.PlayerTypeCount(gameId, playerType)
	if err != nil {
		return err
	}
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
