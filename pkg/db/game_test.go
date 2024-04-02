package db_test

import (
	"fmt"
	"testing"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/google/uuid"
)

func TestGetAllHands(t *testing.T) {
	db.SetConfig(&db.Config{
		Addr: "localhost:6379",
	})

	gameId, _ := uuid.Parse("018e9d85-233b-72c5-b14c-c7ba8699616e")

	val, err := db.GetMultiPlayerProperty[[]int](gameId, "hand")
	if err != nil {
		t.Error(err)
		return
	}

	all := []int{}

	for _, v := range val {
		all = append(all, v...)
	}

	fmt.Printf("%v\n", all)
}
