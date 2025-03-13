package wal

import (
	"errors"
	"io"
	"os"
)

type Buffer interface {
	WriteAt(data []byte, offset int64) (n int, err error)
	ReadAt(data []byte, offset int64) (n int, err error)
	io.Reader
	Seek(offset int64, whence int) (ret int64, err error)
}

func OpenFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
}

func NewBytesBuffer(size int) *bytesBuffer {
	return &bytesBuffer{buf: make([]byte, 0, size)}
}

// bytesBuffer 是一个实现 Buffer 接口的结构体，用于存储字节数据。
type bytesBuffer struct {
	buf []byte

	r int64
	w int64
}

func (b *bytesBuffer) Seek(offset int64, whence int) (ret int64, err error) {
	b.r = offset
	b.w = offset
	return offset, nil
}

func (b *bytesBuffer) Read(p []byte) (n int, err error) {
	n, err = b.ReadAt(p, b.r)
	if err != nil {
		return
	}
	b.r += int64(n)
	return
}

func (b *bytesBuffer) Write(p []byte) (n int, err error) {
	n, err = b.WriteAt(p, b.w)
	if err != nil {
		return
	}
	b.w += int64(n)
	return
}

func (b *bytesBuffer) WriteAt(data []byte, offset int64) (int, error) {
	if offset < 0 {
		return 0, errors.New("invalid offset")
	}
	totalLength := offset + int64(len(data))
	if totalLength > int64(len(b.buf)) {
		newBuf := make([]byte, totalLength)
		copy(newBuf, b.buf)
		b.buf = newBuf
	}

	copy(b.buf[offset:], data)
	return len(data), nil
}

func (b *bytesBuffer) ReadAt(data []byte, offset int64) (n int, err error) {
	if offset < 0 {
		return 0, errors.New("invalid offset")
	}
	if offset >= int64(len(b.buf)) {
		return 0, io.EOF
	}
	n = len(b.buf) - int(offset)
	if n > len(data) {
		n = len(data)
	}
	copy(data[:n], b.buf[offset:])
	return n, nil
}
