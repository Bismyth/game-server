package skull

import (
	"fmt"

	"github.com/Bismyth/game-server/db"
	"github.com/Bismyth/game-server/interfaces"
	"github.com/google/uuid"
)

func cleanup(gameId uuid.UUID) error {
	err := C_PLAYER.For(gameId).Delete()

	err = db.ExpireCache(gameId, cacheExpireTime)
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

func getTileList(gameId uuid.UUID) ([]int, error) {
	players, err := C_PLAYER.For(gameId).GetAll()
	if err != nil {
		return []int{}, err
	}

	placedTiles := make([]int, len(players))
	for i, playerId := range players {
		t := PD_TILES_PLACED.MustGetPD(gameId, playerId)
		placedTiles[i] = len(t)
	}

	return placedTiles, nil
}

func countTilesPlaced(gameId uuid.UUID) ([]int, error) {
	players, err := C_PLAYER.For(gameId).GetAll()
	if err != nil {
		return []int{}, err
	}

	placedTiles := make([]int, len(players))
	for i, playerId := range players {
		t := PD_TILES_PLACED.MustGetPD(gameId, playerId)
		placedTiles[i] = len(t)
	}

	return placedTiles, nil
}

func totalTilesPlaced(gameId uuid.UUID) (int, error) {
	tiles, err := getTileList(gameId)
	if err != nil {
		return 0, err
	}
	sum := 0
	for _, num := range tiles {
		sum += num
	}
	return sum, nil
}

func allPlayersTilesPlaced(gameId uuid.UUID) (bool, error) {
	tiles, err := getTileList(gameId)
	if err != nil {
		return false, err
	}

	output := true
	for _, num := range tiles {
		if num <= 0 {
			output = false
		}
	}
	return output, nil
}

func goToUnpassedPlayer(gameId uuid.UUID, passed []uuid.UUID) (uuid.UUID, error) {
	cursor := C_PLAYER.For(gameId)
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
	gs, err := cachePublicGameState(gameId)
	if err != nil {
		return err
	}

	c.SendGlobal(GameState{Public: gs})

	return nil
}

func newRound(c interfaces.GameCommunication, gameId uuid.UUID) error {
	roundNumber := PGS_ROUND.MustGet(gameId)
	PGS_ROUND.MustSet(gameId, roundNumber+1)
	resetRoundValues(gameId)

	err := updatePublicGameState(c, gameId)
	if err != nil {
		return err
	}

	players, err := C_PLAYER.For(gameId).GetAll()
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

	return nil
}

func resetRoundValues(gameId uuid.UUID) error {
	PGS_BID.MustSet(gameId, 0)
	PGS_FLIPPER.MustSet(gameId, uuid.Nil)
	PGS_PASSED.MustSet(gameId, []uuid.UUID{})
	PGS_PLAYER_LEFT.MustSet(gameId, false)

	players, err := C_PLAYER.For(gameId).GetAll()
	if err != nil {
		return err
	}
	for _, player := range players {
		PD_TILES_PLACED.MustSet(gameId, player, []Tile{})
		PD_TILES_REVEALED.MustSet(gameId, player, 0)
	}

	return nil
}

func endGame(c interfaces.GameCommunication, gameId uuid.UUID) error {
	PGS_GAME_OVER.MustSet(gameId, true)

	pgs, err := cachePublicGameState(gameId)
	if err != nil {
		return err
	}

	c.EndGame()

	c.SendGlobal(GameState{
		Public: pgs,
	})

	err = cleanup(gameId)
	if err != nil {
		return err
	}

	return nil
}

func checkEnd(gameId uuid.UUID) (bool, error) {
	numPlayers, err := C_PLAYER.For(gameId).Length()
	if err != nil {
		return false, err
	}

	return numPlayers <= 1, nil
}
