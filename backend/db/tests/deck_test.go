package db

import (
	"fmt"
	"testing"

	"github.com/Bismyth/game-server/db"
	"github.com/google/uuid"
)

func TestCursor(t *testing.T) {

	db.SetConfig(db.Config{
		Addr: "localhost:6379",
	})

	gameId, err := uuid.NewV7()
	if err != nil {
		t.Fatal(err)
	}

	DECK := db.Deck[int]{Key: "deck"}

	d := DECK.For(gameId)

	for i := 0; i < 10; i++ {
		d.Add(i)
	}

	fmt.Printf("%v\n", d.GetAll())

	d.Shuffle()

	fmt.Printf("%v\n", d.GetAll())
}
