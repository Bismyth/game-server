package server

import (
	"log"
	"net/http"
	"os"
	"path"

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
		http.Handle("/", ServeStaticSPA(S.Config.Application.PublicDir))
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

func ServeStaticSPA(servePath string) http.Handler {
	fs := http.Dir(servePath)
	fsh := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fs.Open(path.Clean(r.URL.Path))
		if os.IsNotExist(err) {
			http.ServeFile(w, r, path.Join(servePath, "index.html"))
			return
		}
		fsh.ServeHTTP(w, r)
	})
}
