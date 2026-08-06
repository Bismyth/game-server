package skull

import (
	"encoding/json"
	"fmt"

	"github.com/Bismyth/game-server/interfaces"
	"github.com/google/uuid"
)

func handlePlace(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID, data json.RawMessage) error {
	var placeData ActionPlace
	err := json.Unmarshal(data, &placeData)
	if err != nil {
		return err
	}
	currentTilesPlaced := PD_TILES_PLACED.MustGetPD(gameId, playerId)

	bid := PGS_BID.MustGet(gameId)
	if bid > 0 {
		return fmt.Errorf("cant place tile if bid has been made")
	}

	hand := PD_TILES.MustGetPD(gameId, playerId)
	if len(currentTilesPlaced) == len(hand) {
		return fmt.Errorf("no more tiles to place")
	}

	inHand := countTiles(hand, placeData.Tile)
	placed := countTiles(currentTilesPlaced, placeData.Tile)
	if inHand-placed <= 0 {
		return fmt.Errorf("no tile left of that kind")
	}

	currentTilesPlaced = append(currentTilesPlaced, placeData.Tile)
	PD_TILES_PLACED.MustSet(gameId, playerId, currentTilesPlaced)

	if len(currentTilesPlaced) > 1 {
		_, err := C_PLAYER.For(gameId).Next()
		if err != nil {
			return err
		}
	}

	err = updatePublicGameState(c, gameId)
	if err != nil {
		return err
	}
	ps, err := getPrivateGameState(gameId, playerId)
	if err != nil {
		return err
	}
	c.SendPlayer(playerId, GameState{Private: ps})

	return nil
}
