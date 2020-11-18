package main

import (
	"github.com/minghex/cacheDB/rodis/cache"
	"github.com/minghex/cacheDB/rodis/http"
)

func main() {
	s := http.NewServer(cache.NewInmemory())
	s.Serve()
}
