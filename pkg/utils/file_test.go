package utils

import (
	"bytes"
	"fmt"
	"math/rand"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFileRelativePath(t *testing.T) {
	getResult := func(v1, v2 string) string {
		path, err := GetFileRelativePath(v1, v2)
		if err != nil {
			t.Fatal(err)
		}
		return path
	}
	assert.Equal(t, getResult("/data/log/test/a.log", "/data"), "log/test/a.log")

	t.Log(filepath.Rel("/data", "/data/log/test/a.log"))
	rel, _ := filepath.Rel("/data", "/data/log/test/a.log")
	t.Log(filepath.Join("/data", rel))
}

func TestGetGoProjectDir(t *testing.T) {
	t.Log(GetGoProjectDir())
	curDir, _ := filepath.Abs(".")
	t.Log(curDir)
}

func TestGetCmdName(t *testing.T) {
	t.Log(GetCmdName())
}

func TestCheckStdInFromPiped(t *testing.T) {
	t.Log(CheckStdInFromPiped())
}

func TestReadLineByFunc(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		result, err := ReadLines(bytes.NewBufferString(`hello
world
!`))
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, result, []string{"hello", "world", "!"})
	})
	t.Run("error", func(t *testing.T) {
		result := make([]string, 0)
		err := ReadLineByFunc(bytes.NewBufferString(`hello
world
!`), func(line string) error {
			if line == "!" {
				return fmt.Errorf(`invalid string, str: %v`, line)
			}
			result = append(result, line)
			return nil
		})
		if err != nil {
			assert.Error(t, err)
		}
		assert.Equal(t, result, []string{"hello", "world"})
	})
	t.Run("prefix", func(t *testing.T) {
		line1 := RandString(rand.Intn(4096))
		line2 := RandString(4096 + rand.Intn(4096))
		line3 := RandString(4096*2 + rand.Intn(4096))
		lines := make([]string, 0)
		if err := ReadLineByFunc(bytes.NewBufferString(line1+"\n"+line2+"\n"+line3+"\n"+"1"+"\n"), func(line string) error {
			lines = append(lines, line)
			return nil
		}); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, lines, []string{line1, line2, line3, "1"})
	})
}
