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

	mux := http.NewServeMux()

	if S.Config.Application.Production {
		mux.Handle("/", ServeStaticSPA(S.Config.Application.PublicDir))
	}

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(S.WSHub, w, r)
	})

	S.RegisterAPI(mux)

	log.Printf("Serving on: %s\n", S.Config.Application.BindAddress)
	err := http.ListenAndServe(S.Config.Application.BindAddress, mux)
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
