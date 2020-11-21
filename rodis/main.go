package main

import (
	"flag"

	"github.com/minghex/cacheDB/rodis/cache"
	"github.com/minghex/cacheDB/rodis/http"
	"github.com/minghex/cacheDB/rodis/tcp"
)

func main() {
	c := cache.NewInmemory()
	go tcp.NewServer(c).Serve()
	http.NewServer(c).Serve()
}

var typ string

func init() {
	flag.StringVar(&typ, "type", "tcp", "server type http/tcp")
}
