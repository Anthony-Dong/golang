package gui

import (
	"bytes"
	"fmt"
	"io"

	"github.com/mattn/go-runewidth"

	"github.com/anthony-dong/golang/pkg/utils"
)

func Printf(w io.Writer, format string, args ...interface{}) {
	if len(args) == 0 {
		w.Write([]byte(terminalLine(format)))
		return
	}
	w.Write([]byte(terminalLine(fmt.Sprintf(format, args...))))
}

func Println(w io.Writer, format string, args ...interface{}) {
	if len(args) == 0 {
		w.Write([]byte(terminalLine(format)))
		w.Write([]byte{'\n'})
		return
	}
	w.Write([]byte(terminalLine(fmt.Sprintf(format, args...))))
	w.Write([]byte{'\n'})
	utils.PrettyJson("")
}

func terminalLine(line string) string {
	var lineWithWidth = bytes.NewBuffer(make([]byte, 0, len(line)))
	for _, r := range line {
		w := runewidth.RuneWidth(r)
		if w == 0 {
			w = 1
		}
		lineWithWidth.WriteString(string(r))
		for i := 1; i < w; i++ {
			lineWithWidth.WriteString(" ")
		}
	}
	return lineWithWidth.String()
}
