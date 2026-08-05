package skull

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"

	"github.com/Bismyth/game-server/pkg/db"
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
	tp, err := GetPlayerProperty[[]Tile](gameId, flipData.Player, pd_tilesPlaced)
	if err != nil {
		return fmt.Errorf("failed to fetch targets placed tiles")
	}

	if len(tr) >= len(tp) {
		return fmt.Errorf("player has had all tiles flipped")
	}

	tilesRevealed := len(tr) + 1
	tile := tp[len(tp)-tilesRevealed]
	err = SetPlayerProperty(gameId, flipData.Player, pd_tilesRevealed, tilesRevealed)
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

	err = updatePublicGameState(c, gameId)
	if err != nil {
		return err
	}

	if len(hand) <= 1 {
		removeActivePlayer(gameId, playerId)
		end, err := checkEnd(gameId)
		if err != nil {
			return err
		}
		if end {
			endGame(c, gameId)
		}
	} else {
		randomIndex := rand.IntN(len(hand))
		hand[randomIndex] = hand[len(hand)-1]
		hand = hand[:len(hand)-1]
		err = SetPlayerProperty(gameId, playerId, pd_tiles, hand)
		if err != nil {
			return err
		}
	}

	err = newRound(c, gameId)
	if err != nil {
		return err
	}

	return nil
}

func removeActivePlayer(gameId uuid.UUID, playerId uuid.UUID) error {
	isPlayer := db.PlayerIsType(gameId, playerId, playerType)
	if !isPlayer {
		return nil
	}
	cursor := db.GetCursor(gameId, playerType)
	current, err := cursor.Current()
	if err != nil {
		return err
	}
	if current == playerId {
		err := cursor.Remove()
		if err != nil {
			return err
		}
	} else {
		err := cursor.SeekIndex(playerId)
		if err != nil {
			return err
		}
		err = cursor.Remove()
		if err != nil {
			return err
		}
		err = cursor.SeekIndex(current)
		if err != nil {
			return err
		}
	}

	return nil
}

func flippedBid(c interfaces.GameCommunication, gameId, playerId uuid.UUID) error {
	err := updatePublicGameState(c, gameId)
	if err != nil {
		return err
	}

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
		err = endGame(c, gameId)
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

	err = SetPlayerProperty(gameId, playerId, pd_tilesRevealed, len(tiles))
	if err != nil {
		return err
	}

	bid, err := GetProperty[int](gameId, d_bid)
	if err != nil {
		return err
	}

	roses := 0

	for i := len(tiles) - 1; i >= 0; i-- {
		if tiles[i] == Skull {
			err = flippedSkull(c, gameId, playerId)
			if err != nil {
				return err
			}
			return nil
		} else {
			roses += 1
			if roses >= bid {
				err = SetPlayerProperty(gameId, playerId, pd_tilesRevealed, len(tiles)-i)
				if err != nil {
					return err
				}
				flippedBid(c, gameId, playerId)
			}
		}
	}

	currentBid, err := GetProperty[int](gameId, d_bid)
	if err != nil {
		return err
	}

	err = updatePublicGameState(c, gameId)
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

	return nil
}
