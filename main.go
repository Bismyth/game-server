package main

import (
	"embed"
	"fmt"
	"log"

	"github.com/Bismyth/game-server/pkg/config"
	"github.com/Bismyth/game-server/pkg/db"
	"github.com/Bismyth/game-server/pkg/server"
	"github.com/Bismyth/game-server/pkg/ws"
)

//go:embed .output/*
var frontendFS embed.FS

func main() {

	c := config.New()

	fmt.Println(c.Redis.Addr)

	err := db.FlushDB()
	if err != nil {
		log.Fatal("could not flush existing db")
	}

	hub := ws.NewHub()
	go hub.Run()

	s := server.New(c, hub, frontendFS)
	s.Run()
}
