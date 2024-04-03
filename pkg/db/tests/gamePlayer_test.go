package db_test

import (
	"testing"

	"github.com/Bismyth/game-server/pkg/db"
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

	playerType := "player"

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
		err := db.PlayerGiveType(gameId, user, playerType)
		if err != nil {
			t.Fatal(err)
		}
	}

	c := db.NewCursor(gameId, playerType)

	current, err := c.Current()
	if err != nil {
		t.Fatal(err)
	}
	if current != users[0] {
		t.Fatalf("wrong uuid expected %q, got %q", users[0].String(), current.String())
	}

	n, err := c.Next()
	if err != nil {
		t.Fatal(err)
	}
	if n != users[1] {
		t.Fatalf("wrong uuid expected %q, got %q", users[1].String(), current.String())
	}

	n, err = c.Next()
	if err != nil {
		t.Fatal(err)
	}
	if n != users[0] {
		t.Fatalf("wrong uuid expected %q, got %q", users[0].String(), current.String())
	}

	n, err = c.Previous()
	if err != nil {
		t.Fatal(err)
	}
	if n != users[1] {
		t.Fatalf("wrong uuid expected %q, got %q", users[1].String(), current.String())
	}

	err = c.Remove()
	if err != nil {
		t.Fatal(err)
	}

	n, err = c.Current()
	if err != nil {
		t.Fatal(err)
	}
	if n != users[0] {
		t.Fatalf("wrong uuid expected %q, got %q", users[0].String(), current.String())
	}

	n, err = c.Next()
	if err != nil {
		t.Fatal(err)
	}
	if n != users[0] {
		t.Fatalf("wrong uuid expected %q, got %q", users[0].String(), current.String())
	}
}
