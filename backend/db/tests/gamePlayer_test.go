package db_test

import (
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

	C_PLAYER := db.Cursor{Key: "player"}

	users := make([]uuid.UUID, 2)
	for i := range users {
		id, err := uuid.NewV7()
		if err != nil {
			t.Fatal(err)
		}
		users[i] = id
	}

	t.Log(users)

	for _, user := range users {
		C_PLAYER.For(gameId).Add(user)
	}

	c := C_PLAYER.For(gameId)
	c.Reset()

	current := c.Current()
	if current != users[0] {
		t.Fatalf("wrong uuid expected %q, got %q", users[0].String(), current.String())
	}

	n := c.Next()
	if n != users[1] {
		t.Fatalf("wrong uuid expected %q, got %q", users[1].String(), current.String())
	}

	n = c.Next()
	if n != users[0] {
		t.Fatalf("wrong uuid expected %q, got %q", users[0].String(), current.String())
	}

	n = c.Previous()
	if n != users[1] {
		t.Fatalf("wrong uuid expected %q, got %q", users[1].String(), current.String())
	}

	c.Remove()

	n = c.Current()
	if n != users[0] {
		t.Fatalf("wrong uuid expected %q, got %q", users[0].String(), current.String())
	}

	n = c.Next()
	if n != users[0] {
		t.Fatalf("wrong uuid expected %q, got %q", users[0].String(), current.String())
	}
}
