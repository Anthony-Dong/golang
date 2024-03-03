package record

import (
	"io"
	"os"
)

type consulStorage struct {
	Format string // format / json
	writer io.Writer
}

func (c *consulStorage) Write(data []byte) error {
	if _, err := c.writer.Write(data); err != nil {
		return err
	}
	if _, err := c.writer.Write([]byte("\n")); err != nil {
		return err
	}
	return nil
}

func (c *consulStorage) Read() ([]byte, error) { return nil, nil }

func (c *consulStorage) Close() error { return nil }

func NewConsulStorage() Storage {
	return &consulStorage{writer: os.Stdout, Format: "format"}
}

func NewJsonConsulStorage() Storage {
	return &consulStorage{writer: os.Stdout, Format: "json"}
}

func isConsulStorage(storage Storage) bool {
	if v, _ := storage.(*consulStorage); v != nil {
		return v.Format == "format"
	}
	return false
}

func isConsulJsonStorage(storage Storage) bool {
	if v, _ := storage.(*consulStorage); v != nil {
		return v.Format == "json"
	}
	return false
}
