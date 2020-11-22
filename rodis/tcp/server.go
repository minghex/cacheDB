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
	r := bufio.NewReader(conn)
	resultChan := make(chan chan *result, 0)
	defer close(resultChan)
	go this.reply(conn, resultChan)

	for {
		op, e := r.ReadByte()
		if e != nil {
			if e != io.EOF {
				log.Println("close connection due to error:", e)
			}
			return
		}

		if op == 'S' {
			this.set(resultChan, r)
		} else if op == 'G' {
			this.get(resultChan, r)
		} else if op == 'D' {
			this.del(resultChan, r)
		} else {
			log.Println("unknown op", op)
		}

		if e != nil {
			log.Println("close connection")
			return
		}
	}
}

func (this *Server) reply(conn net.Conn, resultCh chan chan *result) {
	defer conn.Close()
	for {
		c, open := <-resultCh
		if !open {
			return
		}

		r := <-c
		e := sendResponse(r.value, r.err, conn)
		if e != nil {
			log.Println("sendResponse error: ", e)
			return
		}
	}
}

func (this *Server) set(resultCh chan chan *result, r *bufio.Reader) {
	ch := make(chan *result)
	resultCh <- ch

	key, value, err := readKeyAndValue(r)
	if err != nil {
		ch <- &result{nil, err}
		return
	}

	go func() {
		ch <- &result{nil, this.Set(key, value)}
	}()
}

func (this *Server) get(resultCh chan chan *result, r *bufio.Reader) {
	ch := make(chan *result)
	resultCh <- ch

	key, err := readKey(r)
	if err != nil {
		ch <- &result{nil, err}
		return
	}

	go func() {
		v, e := this.Get(key)
		ch <- &result{v, e}
	}()
}

func (this *Server) del(resultCh chan chan *result, r *bufio.Reader) {
	ch := make(chan *result)
	resultCh <- ch
	key, err := readKey(r)
	if err != nil {
		ch <- &result{nil, err}
	}
	go func() {
		ch <- &result{nil, this.Del(key)}
	}()
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
