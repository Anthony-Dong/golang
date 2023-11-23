package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

func PrettyCmd(cmd *exec.Cmd, max ...int) string {
	maxSize := 128
	if len(max) > 0 {
		maxSize = max[0]
	}
	environ := os.Environ()
	writer := NewCommandPretty(maxSize)
	for _, elem := range cmd.Env {
		if Contains(environ, elem) {
			continue
		}
		writer.Write(elem)
	}
	for _, elem := range cmd.Args {
		writer.Write(elem)
	}
	return writer.String()
}

type CommandPretty struct {
	num    int
	max    int
	output strings.Builder
}

func NewCommandPretty(max int) *CommandPretty {
	return &CommandPretty{
		max: max,
	}
}

func (m *CommandPretty) Write(arg string) {
	if m.num > m.max {
		m.output.WriteString("\\\n")
		m.num = 0
	}
	arg = unquoteArg(arg)
	m.output.WriteString(arg)
	m.output.WriteString(" ")
	m.num = m.num + len(arg) + 1
}

func (m *CommandPretty) NewLine() {
	if m.num == 0 {
		return
	}
	m.output.WriteString("\\\n")
	m.num = 0
}

func (m *CommandPretty) String() string {
	return m.output.String()
}

func unquoteArg(arg string) string {
	for _, elem := range arg {
		if unicode.IsSpace(elem) {
			return strconv.Quote(arg)
		}
	}
	return arg
}

func PrettyDir(dir string) string {
	if filepath.IsAbs(dir) {
		return dir
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return dir
	}
	getwd, err := os.Getwd()
	if err != nil {
		return dir
	}

	rel, err := filepath.Rel(getwd, abs)
	if err != nil {
		return dir
	}
	return rel
}
