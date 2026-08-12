package nothanks

import (
	"slices"

	"github.com/Bismyth/game-server/db"
	"github.com/Bismyth/game-server/interfaces"
	"github.com/google/uuid"
)

var TOKEN_AMOUNT = map[int]int{
	2: 11,
	3: 11,
	4: 11,
	5: 11,
	6: 9,
	7: 7,
}

func newRound(gameId uuid.UUID) {
	deck := DECK.For(gameId)

	deck.Clear()

	for x := 3; x < 36; x++ {
		deck.Add(x)
	}

	deck.Shuffle()
	removed := make([]int, 9)

	for x := 0; x < 9; x++ {
		v := deck.Draw()
		removed[x] = v
	}
	GD_REMOVED.MustSet(gameId, removed)

	playerCount := C_PLAYER.For(gameId).Length()

	tokenCount := TOKEN_AMOUNT[int(playerCount)]

	players := C_PLAYER.For(gameId).GetAll()
	for _, player := range players {
		PD_TOKENS.MustSet(gameId, player, tokenCount)
		PD_CARDS.MustSet(gameId, player, []int{})
	}

	first := deck.Draw()
	GD_INPLAY.MustSet(gameId, first)
	GD_TOKENS_ON_CARD.MustSet(gameId, 0)

}

func scorePlayer(gameId, playerId uuid.UUID) int {
	cards := PD_CARDS.MustGet(gameId, playerId)
	tokens := PD_TOKENS.MustGet(gameId, playerId)

	score := 0
	prev := 0
	slices.Sort(cards)
	for _, c := range cards {
		if prev != c-1 {
			score += c
		}
		prev = c
	}
	score -= tokens
	return score
}

func endRound(c interfaces.GameCommunication, gameId uuid.UUID) error {
	players := C_PLAYER.For(gameId).GetAll()

	roundScore := make(map[uuid.UUID]int)
	for _, player := range players {
		score := scorePlayer(gameId, player)
		currentScore := PD_SCORE.MustGet(gameId, player)
		PD_SCORE.MustSet(gameId, player, currentScore+score)
		roundScore[player] = score
	}

	round := GD_ROUND.MustGet(gameId)
	totalRounds := GD_TOTAL_ROUNDS.MustGet(gameId)

	hands, err := PD_CARDS.GetMap(gameId, players)
	if err != nil {
		panic("failed to get player hands")
	}
	removed := GD_REMOVED.MustGet(gameId)

	pr := PreviousRound{
		Type:        "previous",
		Score:       roundScore,
		PlayerCards: hands,
		Removed:     removed,
		Round:       round,
	}

	cachePrevious(gameId, &pr)

	if round == totalRounds {
		return endGame(c, gameId, &pr)
	}

	GD_ROUND.MustSet(gameId, round+1)
	newRound(gameId)

	nextPlayer := C_PLAYER.For(gameId).Current()
	ps := cachePublicGameState(gameId)
	c.SendGlobal(&pr)
	for _, player := range players {
		p := loadPrivate(gameId, player)
		c.SendPlayer(player, NewGameState(ps, &p))
	}

	c.ActionPrompt(nextPlayer, allActions)

	return nil
}

func endGame(c interfaces.GameCommunication, gameId uuid.UUID, pr *PreviousRound) error {
	err := GD_GAME_OVER.Set(gameId, true)
	if err != nil {
		return err
	}

	pGs := cachePublicGameState(gameId)
	c.EndGame()
	c.SendGlobal(NewGameState(pGs, nil))
	if pr != nil {
		c.SendGlobal(pr)
	}

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

	DECK.For(gameId).Clear()

	err = db.ExpireCache(gameId, cacheExpireTime)
	if err != nil {
		return err
	}
	err = db.ExpireGameKey(gameId, "pr", cacheExpireTime)
	if err != nil {
		return err
	}

	err = db.DeleteGame(gameId)
	if err != nil {
		return err
	}

	return nil
}
