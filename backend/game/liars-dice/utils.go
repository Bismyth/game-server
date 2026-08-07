package liarsdice

import (
	"fmt"
	"math/rand/v2"

	"github.com/Bismyth/game-server/db"
	"github.com/Bismyth/game-server/db/msg"
	"github.com/Bismyth/game-server/interfaces"
	"github.com/google/uuid"
)

func progressTurn(c interfaces.GameCommunication, gameId uuid.UUID) error {
	nextPlayer, err := C_PLAYER.For(gameId).Next()
	if err != nil {
		return err
	}

	publicGs, err := cachePublicGameState(gameId)
	if err != nil {
		return err
	}

	gs := GameState{Public: publicGs}

	c.SendGlobal(gs)
	c.ActionPrompt(nextPlayer, allActions)

	return nil
}

func newRound(c interfaces.GameCommunication, gameId uuid.UUID, pr *RoundInfo) error {
	players, err := C_PLAYER.For(gameId).GetAll()
	if err != nil {
		return err
	}

	err = GD_PREVIOUS_ROUND.Set(gameId, *pr)
	if err != nil {
		return err
	}

	err = rollHands(gameId, players)
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

	err = GD_BID.Set(gameId, "")
	if err != nil {
		return err
	}

	err = progressTurn(c, gameId)
	if err != nil {
		return err
	}

	return nil
}

func rollHands(gameId uuid.UUID, players []uuid.UUID) error {
	for _, playerId := range players {
		numDice, err := PD_DICE.Get(gameId, playerId)
		if err != nil {
			return err
		}

		hand := make([]int, numDice)
		for i := range hand {
			hand[i] = rand.IntN(6) + 1
		}

		err = PD_HAND.Set(gameId, playerId, hand)
		if err != nil {
			return err
		}
	}

	return nil
}

func generatePreviousRound(gameId uuid.UUID, pvInfo *ParsedRoundInfo) (*RoundInfo, error) {
	r, err := GD_PREVIOUS_ROUND.Get(gameId)
	if err != nil {
		return nil, err
	}

	var roundInfo RoundInfo

	roundInfo.Round = r.Round + 1

	if pvInfo.Leave != "" {
		roundInfo.Leave = pvInfo.Leave
		return &roundInfo, nil
	}

	players, err := C_PLAYER.For(gameId).GetAll()
	if err != nil {
		return nil, err
	}

	hb, err := GD_BID.Get(gameId)
	if err != nil {
		return nil, err
	}
	roundInfo.HighestBid = hb

	roundInfo.Hands = make(map[uuid.UUID][]int)

	for _, id := range players {
		h, err := PD_HAND.Get(gameId, id)
		if err != nil {
			return nil, err
		}
		roundInfo.Hands[id] = h
	}

	if pvInfo != nil {
		roundInfo.CallUser = pvInfo.CallUser
		roundInfo.DiceLost = pvInfo.DiceLost
		roundInfo.LastBid = pvInfo.LastBid
	}

	return &roundInfo, nil
}

func endGame(c interfaces.GameCommunication, gameId uuid.UUID, pr *RoundInfo) error {
	err := GD_GAME_OVER.Set(gameId, true)
	if err != nil {
		return err
	}

	err = GD_PREVIOUS_ROUND.Set(gameId, *pr)
	if err != nil {
		return err
	}

	pGs, err := cachePublicGameState(gameId)
	if err != nil {
		return err
	}

	c.EndGame()
	c.SendGlobal(GameState{
		Public: pGs,
	})

	err = cleanup(gameId)
	if err != nil {
		return err
	}

	return nil
}

func cleanup(gameId uuid.UUID) error {
	err := C_PLAYER.For(gameId).Delete()
	if err != nil {
		return err
	}

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

func bidMsg(bid string) msg.Part {
	n, v, err := parseBid(bid)
	if err != nil {
		panic(err)
	}

	return msg.Concat(msg.Text(fmt.Sprintf("%dx ", n)), msg.Icon(fmt.Sprintf("bi:dice-%d-fill", v)))
}
