package http

import (
	"net/http"

	"github.com/minghex/cacheDB/rodis/cache"
)

type Server struct {
	//缓存数据的Cache
	cache.Cache
}

func NewServer(c cache.Cache) *Server {
	return &Server{c}
}

func (s *Server) Serve() {
	http.Handle("/cache/", s.cacheHandler())
	http.Handle("/status", s.statusHandler())
	http.ListenAndServe(":12345", nil)
}
