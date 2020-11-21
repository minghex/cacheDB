package tcp

import (
	"bufio"
	"fmt"
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

		if e != nil {
			log.Println("close connection")
			return
		}
	}
}

func (this *Server) set(conn net.Conn, r *bufio.Reader) error {
	key, value, err := readKeyAndValue(r, conn)
	if err != nil {
		return err
	}
	return sendResponse(nil, this.Set(key, value), conn)
}

func (this *Server) get(conn net.Conn, r *bufio.Reader) error {
	key, err := readKey(r)
	if err != nil {
		return err
	}

	value, err := this.Get(key)
	return sendResponse(value, err, conn)
}

func (this *Server) del(conn net.Conn, r *bufio.Reader) error {
	key, err := readKey(r)
	if err != nil {
		return err
	}
	return sendResponse(nil, this.Del(key), conn)
}

func sendResponse(value []byte, err error, conn net.Conn) error {
	if err != nil {
		tmp := fmt.Sprintf("-%d ", len(err.Error())) + err.Error()
		_, e := conn.Write([]byte(tmp))
		return e
	}
	tmp := fmt.Sprintf("%d ", len(value))
	_, e := conn.Write(append([]byte(tmp), value...))
	return e
}
