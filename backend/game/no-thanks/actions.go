package nothanks

import (
	"fmt"

	"github.com/Bismyth/game-server/db"
	"github.com/Bismyth/game-server/db/msg"
	"github.com/Bismyth/game-server/interfaces"
	"github.com/google/uuid"
)

func handlePass(c interfaces.GameCommunication, gameId, playerId uuid.UUID) error {
	currentTokens := PD_TOKENS.MustGet(gameId, playerId)
	if currentTokens <= 0 {
		return fmt.Errorf("not enough tokens")
	}

	db.GameEvent(gameId, c).Log(msg.Msg().Player(playerId).Text(" has passed").String())

	C_PLAYER.For(gameId).Next()
	PD_TOKENS.MustSet(gameId, playerId, currentTokens-1)
	c_t := GD_TOKENS_ON_CARD.MustGet(gameId)
	GD_TOKENS_ON_CARD.MustSet(gameId, c_t+1)

	actionCleanup(c, gameId, playerId)
	return nil
}

func handleTake(c interfaces.GameCommunication, gameId, playerId uuid.UUID) error {
	currentCard := GD_INPLAY.MustGet(gameId)
	hand := PD_CARDS.MustGet(gameId, playerId)

	hand = append(hand, currentCard)
	PD_CARDS.MustSet(gameId, playerId, hand)

	db.GameEvent(gameId, c).Log(msg.Msg().Player(playerId).Text(" has taken ").Bold(msg.Int(currentCard)).String())

	tokensOnCard := GD_TOKENS_ON_CARD.MustGet(gameId)
	playerTokens := PD_TOKENS.MustGet(gameId, playerId)
	PD_TOKENS.MustSet(gameId, playerId, playerTokens+tokensOnCard)

	deck := DECK.For(gameId)

	if deck.Length() <= 0 {
		err := endRound(c, gameId)
		if err != nil {
			return err
		}
		return nil
	}

	newCard := deck.Draw()
	GD_TOKENS_ON_CARD.MustSet(gameId, 0)
	GD_INPLAY.MustSet(gameId, newCard)

	actionCleanup(c, gameId, playerId)

	return nil
}

func actionCleanup(c interfaces.GameCommunication, gameId, playerId uuid.UUID) {
	nextPlayer := C_PLAYER.For(gameId).Current()
	ps := cachePublicGameState(gameId)
	p := loadPrivate(gameId, playerId)

	c.SendGlobal(NewGameState(ps, nil))
	c.SendPlayer(playerId, NewGameState(nil, &p))
	c.ActionPrompt(nextPlayer, allActions)
}
