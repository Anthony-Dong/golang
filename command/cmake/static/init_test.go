package static

import (
	"testing"

	"github.com/anthony-dong/golang/pkg/utils"
)

func Test_ReadAllFiles(t *testing.T) {
	file, err := ReadAllFiles()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(utils.ToJson(file, true))
}
