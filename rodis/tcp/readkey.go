package tcp

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func readKey(r *bufio.Reader) (string, error) {
	l, e := readLen(r)
	if e != nil {
		return "", e
	}

	buff := make([]byte, l)
	_, e = io.ReadFull(r, buff)
	if e != nil {
		return "", e
	}

	return string(buff), nil
}

func readKeyAndValue(r *bufio.Reader, conn net.Conn) (string, []byte, error) {
	klen, e := readLen(r)
	if e != nil {
		return "", nil, fmt.Errorf("read key length error")
	}

	vlen, e := readLen(r)
	if e != nil {
		return "", nil, fmt.Errorf("read value length error")
	}

	kbuf := make([]byte, klen)
	_, e = io.ReadFull(r, kbuf)
	if e != nil {
		return "", nil, fmt.Errorf("read key error")
	}

	vbuf := make([]byte, vlen)
	_, e = io.ReadFull(r, vbuf)
	if e != nil {
		return "", nil, fmt.Errorf("read value error")
	}

	return string(kbuf), vbuf, nil
}
