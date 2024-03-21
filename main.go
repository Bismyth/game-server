package main

import (
	"github.com/Bismyth/game-server/pkg/db"
	"github.com/Bismyth/game-server/pkg/ws"
	"github.com/gin-gonic/gin"
)

func main() {

	db.SetConfig(&db.Config{
		Addr: "localhost:6379",
		DB:   0,
	})

	router := gin.Default()

	hub := ws.NewHub()
	go hub.Run()

	router.Static("/assets", "./public/assets")
	router.StaticFile("/", "./public/index.html")

	router.GET("/ws", func(c *gin.Context) {
		ws.ServeWs(hub, c.Writer, c.Request)
	})
	router.Run(":8080")
}
