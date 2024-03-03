package proxy

import (
	"bufio"
	"bytes"
)

type copyReader struct {
	buffer bytes.Buffer
	reader *bufio.Reader
}

func (r *copyReader) Read(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)
	if n > 0 {
		r.buffer.Write(p[:n])
	}
	return n, err
}

func (r *copyReader) Peek(n int) ([]byte, error) {
	return r.reader.Peek(n)
}
