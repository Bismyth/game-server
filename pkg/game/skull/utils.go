package skull

import (
	"fmt"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/Bismyth/game-server/pkg/interfaces"
	"github.com/google/uuid"
)

func cleanup(gameId uuid.UUID) error {
	c := db.GetCursor(gameId, playerType)
	err := c.Delete()
	if err != nil {
		return err
	}

	err = db.ExpireCache(gameId, cacheExpireTime)
	if err != nil {
		return err
	}

	err = db.DeletePlayerTypeList(gameId, playerType)
	if err != nil {
		return err
	}

	err = db.DeleteGame(gameId)
	if err != nil {
		return err
	}

	return nil
}

func countTiles(arr []Tile, t Tile) int {
	num := 0
	for _, x := range arr {
		if x == t {
			num++
		}
	}
	return num
}

func countTilesPlaced(gameId uuid.UUID) ([]int, error) {
	players, err := db.PlayerTypeGetAll(gameId, playerType)
	if err != nil {
		return []int{}, err
	}

	placedTiles := make([]int, len(players))
	for i, playerId := range players {
		t, err := GetPlayerProperty[[]Tile](gameId, playerId, pd_tilesPlaced)
		if err != nil {
			return placedTiles, err
		}
		placedTiles[i] = len(t)
	}

	return placedTiles, nil
}

func totalTilesPlaced(placed []int) int {
	sum := 0
	for _, num := range placed {
		sum += num
	}
	return sum
}

func allPlayersTilesPlaced(placed []int) bool {
	output := true
	for _, num := range placed {
		if num <= 0 {
			output = false
		}
	}
	return output
}

func findNextUnpassedPlayer(gameId uuid.UUID, passed []uuid.UUID) (uuid.UUID, error) {
	cursor := db.GetCursor(gameId, playerType)
	currentPlayer, err := cursor.Current()
	if err != nil {
		return uuid.Nil, err
	}

	nextPlayer, err := cursor.Next()
	if err != nil {
		return uuid.Nil, err
	}

	for isPassedPlayer(nextPlayer, passed) {
		nextPlayer, err = cursor.Next()
		if err != nil {
			return uuid.Nil, err
		}
		if nextPlayer == currentPlayer {
			return uuid.Nil, fmt.Errorf("no player found")
		}
	}

	return nextPlayer, nil
}

func isPassedPlayer(playerId uuid.UUID, passed []uuid.UUID) bool {
	for _, player := range passed {
		if playerId == player {
			return true
		}
	}
	return false
}

func updatePublicGameState(c interfaces.GameCommunication, gameId uuid.UUID) error {
	err := cachePublicGameState(gameId)
	if err != nil {
		return err
	}

	gs, err := getPublicGameState(gameId)
	if err != nil {
		return err
	}

	c.SendGlobal(GameState{Public: gs})

	return nil
}

func newRound(c interfaces.GameCommunication, gameId uuid.UUID) error {
	err := SetProperty(gameId, d_bid, 0)
	if err != nil {
		return err
	}
	err = resetRoundValues(gameId)
	if err != nil {
		return err
	}

	return nil
}

func resetRoundValues(gameId uuid.UUID) error {
	err := SetProperty(gameId, d_bid, 0)
	if err != nil {
		return err
	}

	err = SetProperty(gameId, d_flipper, uuid.Nil)
	if err != nil {
		return err
	}

	err = SetProperty(gameId, d_currentTurn, uuid.Nil)
	if err != nil {
		return err
	}

	err = SetProperty(gameId, d_passed, []uuid.UUID{})
	if err != nil {
		return err
	}

	players, err := db.PlayerTypeGetAll(gameId, playerType)
	if err != nil {
		return err
	}
	for _, player := range players {
		err = SetPlayerProperty(gameId, player, pd_tilesPlaced, []Tile{})
		if err != nil {
			return err
		}
		err = SetPlayerProperty(gameId, player, pd_tilesRevealed, 0)
		if err != nil {
			return err
		}
	}

	return nil
}

func endGame(gameId uuid.UUID) error {
	return fmt.Errorf("not implemented")
}
