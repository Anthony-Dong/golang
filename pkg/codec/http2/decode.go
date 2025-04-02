package http2

import (
	"bufio"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
)

func DecodeH2CMessage(reader *bufio.Reader) (*http.Request, error) {
	peek, err := reader.Peek(3)
	if err != nil {
		return nil, err
	}
	if string(peek) != "PRI" {
		return nil, fmt.Errorf("invalid http2 message")
	}
	r, err := http.ReadRequest(reader)
	if err != nil {
		return nil, err
	}
	// 24个字节
	if r.Method == "PRI" && len(r.Header) == 0 && r.URL.Path == "*" && r.Proto == "HTTP/2.0" {
		const expectedBody = "SM\r\n\r\n"
		buf := make([]byte, len(expectedBody))
		n, err := io.ReadFull(reader, buf)
		if err != nil {
			return nil, fmt.Errorf("h2c: error reading client preface: %s", err)
		}
		if string(buf[:n]) == expectedBody {
			return r, nil
		}
	}
	return nil, fmt.Errorf(`invalid h2c request`)
}

func DecodeFrame(reader *bufio.Reader) (http2.Frame, error) {
	framer := http2.NewFramer(nil, reader)
	framer.ReadMetaHeaders = hpack.NewDecoder(1024*16, nil)
	return framer.ReadFrame()
}
