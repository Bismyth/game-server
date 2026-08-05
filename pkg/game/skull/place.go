package skull

import (
	"encoding/json"
	"fmt"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/Bismyth/game-server/pkg/interfaces"
	"github.com/google/uuid"
)

func handlePlace(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID, data json.RawMessage) error {
	var placeData ActionPlace
	err := json.Unmarshal(data, &placeData)
	if err != nil {
		return err
	}

	currentTilesPlaced, err := GetPlayerProperty[[]Tile](gameId, playerId, pd_tilesPlaced)
	if err != nil {
		return err
	}

	turn, err := GetProperty[uuid.UUID](gameId, d_currentTurn)
	if err != nil {
		return err
	}

	if len(currentTilesPlaced) > 0 && turn != playerId {
		return fmt.Errorf("not your turn")
	}

	bid, err := GetProperty[int](gameId, d_bid)
	if err != nil {
		return err
	}

	if bid > 0 {
		return fmt.Errorf("cant place tile if bid has been made")
	}

	hand, err := GetPlayerProperty[[]Tile](gameId, playerId, pd_tiles)
	if err != nil {
		return err
	}

	if len(currentTilesPlaced) == len(hand) {
		return fmt.Errorf("no more tiles to place")
	}

	inHand := countTiles(hand, placeData.Tile)
	placed := countTiles(currentTilesPlaced, placeData.Tile)
	if inHand-placed <= 0 {
		return fmt.Errorf("no tile left of that kind")
	}

	currentTilesPlaced = append(currentTilesPlaced, placeData.Tile)
	err = SetPlayerProperty(gameId, playerId, pd_tilesPlaced, currentTilesPlaced)
	if err != nil {
		return err
	}
	cursor := db.GetCursor(gameId, playerType)
	if turn == uuid.Nil {
		numTilesPlaced, err := countTilesPlaced(gameId)
		if err != nil {
			return err
		}
		if allPlayersTilesPlaced(numTilesPlaced) {
			currentPlayer, err := cursor.Current()
			if err != nil {
				return err
			}
			err = SetProperty(gameId, d_currentTurn, currentPlayer)
			if err != nil {
				return err
			}
		}
	} else {
		nextPlayer, err := cursor.Next()
		if err != nil {
			return err
		}
		err = SetProperty(gameId, d_currentTurn, nextPlayer)
		if err != nil {
			return err
		}
	}

	err = cachePublicGameState(gameId)
	if err != nil {
		return err
	}

	gs, err := getPublicGameState(gameId)
	if err != nil {
		return err
	}

	c.SendGlobal(GameState{Public: gs})

	ps, err := getPrivateGameState(gameId, playerId)
	if err != nil {
		return err
	}
	c.SendPlayer(playerId, GameState{Private: ps})

	return nil
}
