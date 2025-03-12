package wal

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/anthony-dong/golang/pkg/compress"
)

func TestNewWal(t *testing.T) {
	buffer := NewBytesBuffer(1024)
	t.Run("test1", func(t *testing.T) {
		wal, err := NewWal(buffer, 1024*10)
		if err != nil {
			t.Fatal(err)
		}
		if err := wal.SetString("k1", "hello world 1", compress.GzipType); err != nil {
			t.Fatal(err)
		}
		if err := wal.SetString("k1", "hello world 2", compress.GzipType); err != nil {
			t.Fatal(err)
		}
		if err := wal.SetString("k3", "hello world 3", compress.SnappyType); err != nil {
			t.Fatal(err)
		}
		t.Run("test1", func(t *testing.T) {
			data, err := wal.GetString("k1")
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, data, "hello world 2")
		})
		t.Run("test2", func(t *testing.T) {
			data, err := wal.GetString("k3")
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, data, "hello world 3")
		})
	})

	if _, err := buffer.Seek(0, io.SeekStart); err != nil {
		panic(err)
	}

	t.Run("test2", func(t *testing.T) {
		wal, err := NewWal(buffer, 1024*10)
		if err != nil {
			t.Fatal(err)
		}
		if err := wal.SetString("k4", "1", compress.SnappyType); err != nil {
			t.Fatal(err)
		}
		if err := wal.SetString("k4", "2", compress.SnappyType); err != nil {
			t.Fatal(err)
		}
	})
}

func Test_NewWalOutOfRange(t *testing.T) {
	buffer := NewBytesBuffer(1024)

	wal, err := NewWal(buffer, 62)
	if err != nil {
		t.Fatal(err)
	}
	if err := wal.SetString("k1", "hello world", compress.NopType); err != nil {
		t.Fatal(err)
	}
	if err := wal.SetString("k1", "hello world", compress.NopType); err != nil {
		t.Fatal(err)
	}
	err = wal.SetString("k1", "hello world", compress.NopType)
	assert.True(t, IsIndexOutOfRangeErr(err))
}

func TestNewWalFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	indexSize := 1024 * 10
	filename := filepath.Join(dir, "index.data")
	t.Log(filename)
	{
		file, err := OpenFile(filename)
		if err != nil {
			t.Fatal(err)
		}
		wal, err := NewWal(file, indexSize)
		if err != nil {
			t.Fatal(err)
		}
		if err := wal.SetString("k1", "v1", compress.SnappyType); err != nil {
			t.Fatal(err)
		}
		if err := wal.SetString("k1", "v2", compress.SnappyType); err != nil {
			t.Fatal(err)
		}
		if err := wal.SetString("k1", "v3", compress.SnappyType); err != nil {
			t.Fatal(err)
		}
		data, err := wal.GetString("k1")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, data, "v3")

		file.Close()
	}

	{
		file, err := OpenFile(filename)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		wal, err := NewWal(file, indexSize)
		if err != nil {
			t.Fatal(err)
		}

		{
			data, err := wal.GetString("k1")
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, data, "v3")
		}

		if err := wal.SetString("k2", "v2", compress.SnappyType); err != nil {
			t.Fatal(err)
		}

		{
			data, err := wal.GetString("k1")
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, data, "v3")
		}

		{
			data, err := wal.GetString("k2")
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, data, "v2")
		}
	}
}
