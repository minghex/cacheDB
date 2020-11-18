package tcp

import (
	"bufio"
	"io"
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
