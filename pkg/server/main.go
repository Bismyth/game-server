package server

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/Bismyth/game-server/pkg/config"
	"github.com/Bismyth/game-server/pkg/ws"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Config   *config.Config
	Frontend embed.FS
	WSHub    *ws.Hub
}

func New(c *config.Config, wshub *ws.Hub, frontend embed.FS) *Server {
	S := &Server{Config: c, Frontend: frontend, WSHub: wshub}
	return S
}

func (S *Server) Run() {
	if S.Config.Application.Production {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})

	if S.Config.Application.Production {
		subFS, err := fs.Sub(S.Frontend, ".output")
		if err != nil {
			panic(err)
		}
		router.NoRoute(gin.WrapH(http.FileServer(http.FS(subFS))))
	}

	router.GET("/ws", func(c *gin.Context) {
		ws.ServeWs(S.WSHub, c.Writer, c.Request)
	})

	router.Run(S.Config.Application.BindAddress)
}
