package record

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_readChunked(t *testing.T) {
	reader := bufio.NewReader(bytes.NewBuffer([]byte("a\r\n1234567890\r\n")))
	chunked, err := readChunked(reader)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(chunked), "1234567890")
	out := bytes.Buffer{}
	if err = writeChunked(&out, []byte("1234567890")); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, out.String(), "a\r\n1234567890\r\n")
}

func Test_fileRecorder(t *testing.T) {
	mkdirTemp, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("dir: ", mkdirTemp)

	{
		recorder, err := NewLocalStorage(filepath.Join(mkdirTemp, "test.log"))
		if err != nil {
			t.Fatal(err)
		}
		for x := 0; x < 10; x++ {
			if err := recorder.Write([]byte(fmt.Sprintf(`hello world + %d`, x))); err != nil {
				t.Fatal(err)
			}
		}
		recorder.Close()
	}

	{
		recorder, err := NewLocalStorage(filepath.Join(mkdirTemp, "test.log"))
		if err != nil {
			t.Fatal(err)
		}
		for {
			read, err := recorder.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				t.Fatal(err)
			}
			t.Log(string(read))
		}
		recorder.Close()
	}
}
