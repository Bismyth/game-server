package skull

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"

	"github.com/Bismyth/game-server/db"
	"github.com/Bismyth/game-server/interfaces"
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
	tp, err := PD_TILES_PLACED.Get(gameId, flipData.Player)
	if err != nil {
		return fmt.Errorf("failed to fetch targets placed tiles")
	}

	if len(tr) >= len(tp) {
		return fmt.Errorf("player has had all tiles flipped")
	}

	tilesRevealed := len(tr) + 1
	tile := tp[len(tp)-tilesRevealed]

	PD_TILES_REVEALED.MustSet(gameId, flipData.Player, tilesRevealed)

	if tile == Skull {
		err = flippedSkull(c, gameId, playerId, flipData.Player)
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

func flippedSkull(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID, target uuid.UUID) error {
	hand := PD_TILES.MustGetPD(gameId, playerId)

	err := updatePublicGameState(c, gameId)
	if err != nil {
		return err
	}

	if len(hand) <= 1 {
		err := db.RemoveFromCursor(gameId, playerId, playerType)
		if err != nil {
			return err
		}
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
		PD_TILES.MustSet(gameId, playerId, hand)
	}

	cursor := db.GetCursor(gameId, playerType)
	err = cursor.SeekIndex(target)
	if err != nil {
		return err
	}

	err = newRound(c, gameId)
	if err != nil {
		return err
	}

	return nil
}

func flippedBid(c interfaces.GameCommunication, gameId, playerId uuid.UUID) error {
	err := updatePublicGameState(c, gameId)
	if err != nil {
		return err
	}
	currentPoints := PD_POINTS.MustGetPD(gameId, playerId)
	newPoints := currentPoints + 1
	PD_POINTS.MustSet(gameId, playerId, newPoints)

	if newPoints >= 2 {
		err = endGame(c, gameId)
		if err != nil {
			return err
		}
		return nil
	}
	cursor := db.GetCursor(gameId, playerType)
	err = cursor.SeekIndex(playerId)
	if err != nil {
		return err
	}
	_, err = cursor.Next()
	if err != nil {
		return err
	}

	err = newRound(c, gameId)
	if err != nil {
		return err
	}

	return nil
}

func startFlipper(c interfaces.GameCommunication, gameId uuid.UUID, playerId uuid.UUID) error {
	PGS_FLIPPER.MustSet(gameId, playerId)

	tiles := PD_TILES_PLACED.MustGetPD(gameId, playerId)
	PD_TILES_REVEALED.MustSet(gameId, playerId, len(tiles))

	bid := PGS_BID.MustGet(gameId)
	roses := 0

	for i := len(tiles) - 1; i >= 0; i-- {
		if tiles[i] == Skull {
			err := flippedSkull(c, gameId, playerId, playerId)
			if err != nil {
				return err
			}
			return nil
		} else {
			roses += 1
			if roses >= bid {
				PD_TILES_REVEALED.MustSet(gameId, playerId, len(tiles)-i)
				return flippedBid(c, gameId, playerId)
			}
		}
	}

	err := updatePublicGameState(c, gameId)
	if err != nil {
		return err
	}

	return nil
}
