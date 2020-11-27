package cacheClient

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type tcpClient struct {
	net.Conn
	r *bufio.Reader
}

func newTcpClient(addr string) Client {
	conn, err := net.Dial("tcp", addr+":13345")
	if err != nil {
		log.Println("dial tcp server fail, error: ", err)
		panic(err)
	}
	r := bufio.NewReader(conn)
	return &tcpClient{conn, r}
}

func readLen(r *bufio.Reader) int {
	tmp, e := r.ReadString(' ')
	if e != nil {
		log.Println(e)
		return 0
	}
	l, e := strconv.Atoi(strings.TrimSpace(tmp))
	if e != nil {
		log.Println(e)
		return 0
	}
	return l
}

func (this *tcpClient) recvResponse() (string, error) {
	vlen := readLen(this.r)
	if vlen == 0 {
		return "", nil
	}

	value := make([]byte, vlen)
	_, e := io.ReadFull(this.r, value)
	if e != nil {
		return "", e
	}
	return string(value), nil
}

func (this *tcpClient) Run(cmd *Cmd) {
	if cmd.OpName == "set" {
		this.set(cmd.Key, cmd.Value)
		_, cmd.Error = this.recvResponse()
		return
	}
	if cmd.OpName == "get" {
		this.get(cmd.Key)
		cmd.Value, cmd.Error = this.recvResponse()
		return
	}
	if cmd.OpName == "del" {
		this.del(cmd.Key)
		_, cmd.Error = this.recvResponse()
		return
	}

	panic("unknown cmd name " + cmd.OpName)
}

func (this *tcpClient) set(k, v string) {
	klen := len(k)
	vlen := len(v)
	req := fmt.Sprintf("S%d %d %s%s", klen, vlen, k, v)
	this.Write([]byte(req))
}

func (this *tcpClient) get(k string) {
	klen := len(k)
	this.Write([]byte(fmt.Sprintf("G%d %s", klen, k)))
}

func (this *tcpClient) del(k string) {
	klen := len(k)
	this.Write([]byte(fmt.Sprintf("D%d %s", klen, k)))
}

func (this *tcpClient) PipelineRun(cmds []*Cmd) {
	for _, c := range cmds {
		if c.OpName == "get" {
			this.get(c.Key)
		} else if c.OpName == "set" {
			this.set(c.Key, c.Value)
		} else if c.OpName == "del" {
			this.del(c.Key)
		}
	}

	for i := 0; i < len(cmds); i++ {
		cmds[i].Value, cmds[i].Error = this.recvResponse()
	}
}
