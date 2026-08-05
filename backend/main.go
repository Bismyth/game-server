package main

import (
	"github.com/Bismyth/game-server/config"
	"github.com/Bismyth/game-server/server"
	"github.com/Bismyth/game-server/ws"
)

// TODO: implement nicer logging library
func main() {

	c := config.New()

	// err := db.FlushDB()
	// if err != nil {
	// 	log.Fatal("could not flush existing db")
	// }

	hub := ws.NewHub()
	go hub.Run()

	s := server.New(c, hub)
	s.Run()
}
