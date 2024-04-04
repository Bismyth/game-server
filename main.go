package main

import (
	"log"

	"github.com/Bismyth/game-server/pkg/config"
	"github.com/Bismyth/game-server/pkg/db"
	"github.com/Bismyth/game-server/pkg/server"
	"github.com/Bismyth/game-server/pkg/ws"
)

// TODO: implement nicer logging library
func main() {

	c := config.New()

	err := db.FlushDB()
	if err != nil {
		log.Fatal("could not flush existing db")
	}

	hub := ws.NewHub()
	go hub.Run()

	s := server.New(c, hub)
	s.Run()
}
