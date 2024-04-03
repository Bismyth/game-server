package server

import (
	"log"
	"net/http"

	"github.com/Bismyth/game-server/pkg/config"
	"github.com/Bismyth/game-server/pkg/ws"
)

type Server struct {
	Config *config.Config
	WSHub  *ws.Hub
}

func New(c *config.Config, wshub *ws.Hub) *Server {
	S := &Server{Config: c, WSHub: wshub}
	return S
}

func (S *Server) Run() {
	if S.Config.Application.Production {
		fs := http.FileServer(http.Dir(S.Config.Application.PublicDir))
		http.Handle("/", fs)
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(S.WSHub, w, r)
	})

	log.Printf("Serving on: %s\n", S.Config.Application.BindAddress)
	err := http.ListenAndServe(S.Config.Application.BindAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}
