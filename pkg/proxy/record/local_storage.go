package record

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

type Storage interface {
	Write(data []byte) error
	Read() ([]byte, error)
	Close() error
}

func NewLocalStorage(file string) (Storage, error) {
	fd, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &localStorage{file: fd, reader: bufio.NewReader(fd)}, nil
}

func isLocalStorage(storage Storage) bool {
	if storage == nil {
		return false
	}
	_, isOK := storage.(*localStorage)
	return isOK
}

type localStorage struct {
	reader *bufio.Reader
	file   *os.File
}

func (f *localStorage) Close() error {
	return f.file.Close()
}

func (f *localStorage) Write(data []byte) error {
	return writeChunked(f.file, data)
}

func (f *localStorage) Read() ([]byte, error) {
	return readChunked(f.reader)
}

func writeChunked(w io.Writer, data []byte) error {
	if _, err := fmt.Fprintf(w, "%x\r\n", len(data)); err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	if _, err := w.Write([]byte("\r\n")); err != nil {
		return err
	}
	return nil
}

func readChunked(reader *bufio.Reader) ([]byte, error) {
	sizeLine, err := reader.ReadSlice('\n')
	if err != nil {
		return nil, err
	}
	if len(sizeLine) <= 2 {
		return nil, err
	}
	size, err := parseHexUint(sizeLine[:len(sizeLine)-2])
	if err != nil {
		return nil, err
	}
	result := make([]byte, size+2)
	if _, err := io.ReadFull(reader, result); err != nil {
		return nil, err
	}
	if string(result[size:]) != "\r\n" {
		return nil, fmt.Errorf("malformed chunked encoding")
	}
	return result[:size], nil
}

func parseHexUint(v []byte) (n uint64, err error) {
	for i, b := range v {
		switch {
		case '0' <= b && b <= '9':
			b = b - '0'
		case 'a' <= b && b <= 'f':
			b = b - 'a' + 10
		case 'A' <= b && b <= 'F':
			b = b - 'A' + 10
		default:
			return 0, errors.New("invalid byte in chunk length")
		}
		if i == 16 {
			return 0, errors.New("http chunk length too large")
		}
		n <<= 4
		n |= uint64(b)
	}
	return
}
