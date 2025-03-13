package wal

import (
	"testing"

	"github.com/anthony-dong/golang/pkg/compress"
)

func TestWal_Search(t *testing.T) {
	wal, err := NewWal(NewBytesBuffer(1024), 1024*10)
	if err != nil {
		t.Fatal(err)
	}
	if err := wal.Set("k1", []byte("v"), compress.NopType, map[string]string{
		"t1": "v1",
	}); err != nil {
		t.Fatal(err)
	}
	if err := wal.Set("k2", []byte("v"), compress.NopType, map[string]string{
		"t1": "v1",
		"t2": "v2",
	}); err != nil {
		t.Fatal(err)
	}
	t.Log(wal.SearchByTags(map[string]string{"t1": "v1"}))
	t.Log(wal.SearchByTags(map[string]string{"t1": "v1", "t2": "v2"}))
}
