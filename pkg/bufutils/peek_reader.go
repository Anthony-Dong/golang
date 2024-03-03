package bufutils

import (
	"io"
	"time"
)

type BufferedWReader interface {
	io.Reader
	io.Writer
	Peek(int) ([]byte, error)
}

type TimeoutReader interface {
	SetReadDeadline(t time.Time) error
}
