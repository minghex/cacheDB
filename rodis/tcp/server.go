package tcp

import (
	"bufio"
	"io"
	"log"
	"net"

	"github.com/minghex/cacheDB/rodis/cache"
)

type Server struct {
	cache.Cache
}

func NewServer(c cache.Cache) *Server {
	return &Server{c}
}

func (this *Server) Serve() {
	ls, e := net.Listen("tcp", ":13345")
	if e != nil {
		panic(e)
	}

	for {
		conn, e := ls.Accept()
		if e != nil {
			continue
		}
		go this.process(conn)
	}
}

func (this *Server) process(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		op, e := r.ReadByte()
		if e != nil {
			if e != io.EOF {
				log.Println("close connection due to error:", e)
			}
			return
		}

		if op == 'S' {
			e = this.set(conn, r)
		} else if op == 'G' {
			e = this.get(conn, r)
		} else if op == 'D' {
			e = this.del(conn, r)
		} else {
			log.Println("unknown op", op)
		}

	}
}

func (this *Server) set(conn net.Conn, r *bufio.Reader) error {
	return nil
}

func (this *Server) get(conn net.Conn, r *bufio.Reader) error {
	return nil
}

func (this *Server) del(conn net.Conn, r *bufio.Reader) error {
	return nil
}
