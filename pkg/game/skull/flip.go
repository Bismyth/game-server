package skull

import (
	"encoding/json"
	"fmt"

	"github.com/Bismyth/game-server/pkg/interfaces"
	"github.com/google/uuid"
)

func handleFlip(c interfaces.GameCommunication, gameId, playerId uuid.UUID, data json.RawMessage) error {
	currentGameState, err := getPublicGameState(gameId)
	if err != nil {
		return err
	}

	if playerId != currentGameState.Flipper {
		return fmt.Errorf("you are not the current flipper")
	}

	var flipData ActionFlip
	err = json.Unmarshal(data, &flipData)
	if err != nil {
		return err
	}

	tr, ok := currentGameState.TilesRevealed[flipData.Player]
	if !ok {
		return fmt.Errorf("invalid target")
	}
	tp, err := GetPlayerProperty[[]Tile](gameId, playerId, pd_tilesPlaced)
	if err != nil {
		return fmt.Errorf("failed to fetch targets placed tiles")
	}

	if len(tr) >= len(tp) {
		return fmt.Errorf("player has had all tiles flipped")
	}
	tile := tp[len(tp)-len(tr)-1]
	tr = append(tr, tile)
	err = SetPlayerProperty(gameId, playerId, pd_tilesRevealed, tr)
	if err != nil {
		return err
	}

	if tile == Skull {
		err = flippedSkull(c, gameId, playerId)
		if err != nil {
			return err
		}
		return nil
	}

	totalRevealed := 0
	for _, current := range currentGameState.TilesRevealed {
		totalRevealed += len(current)
	}

	if (totalRevealed + 1) >= currentGameState.Bid {
		err = flippedBid(c, gameId, playerId)
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

func flippedSkull(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID) error {
	
	hand, err := GetPlayerProperty[[]Tile](gameId, playerId, pd_tiles)
	if err != nil {
		return err
	}

	if len(hand) <= 1 {
		// make player out
	} else {

	}


	return nil
}

func flippedBid(c interfaces.GameCommunication, gameId, playerId uuid.UUID) error {
	currentPoints, err := GetPlayerProperty[int](gameId, playerId, pd_points)
	if err != nil {
		return err
	}

	newPoints := currentPoints + 1
	err = SetPlayerProperty(gameId, playerId, pd_points, newPoints)
	if err != nil {
		return err
	}

	if newPoints >= 2 {
		err = endGame(gameId)
		if err != nil {
			return err
		}
		return nil
	}

	err = newRound(c, gameId)
	if err != nil {
		return err
	}

	return nil
}

func startFlipper(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID) error {
	err := SetProperty(gameId, d_flipper, playerId)
	if err != nil {
		return err
	}

	tiles, err := GetPlayerProperty[[]Tile](gameId, playerId, pd_tilesPlaced)
	if err != nil {
		return err
	}

	roses := 0
	for _, t := range tiles {
		if t == Skull {
			err = flippedSkull(c, gameId, playerId)
			if err != nil {
				return err
			}
			return nil
		} else {
			roses += 1
		}
	}

	currentBid, err := GetProperty[int](gameId, d_bid)
	if err != nil {
		return err
	}

	if currentBid <= roses {
		err = flippedBid(c, gameId, playerId)
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
