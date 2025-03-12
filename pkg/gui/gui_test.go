package gui

import (
	"testing"

	"github.com/mattn/go-runewidth"
)

func TestName(t *testing.T) {
	data := "你好世界, hello world!"
	for _, elem := range data {
		t.Log(string(elem), ":", runewidth.RuneWidth(elem))
	}
	t.Log(terminalLine(data))
}
