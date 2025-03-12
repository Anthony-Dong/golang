package utils

import (
	"fmt"
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

var unitMap = []string{"Byte", "KB", "MB", "GB", "TB", "PB", "EB"}

// PrettySize 函数用于将字节大小转换为易读的格式
func PrettySize(size int) string {
	// 处理输入为 0 的情况
	if size == 0 {
		return fmt.Sprintf("%d%s", size, unitMap[0])
	}
	// 用于记录转换的次数
	unitIndex := 0
	// 存储最终转换后的大小
	floatSize := float64(size)
	maxUnixIndex := len(unitMap) - 1
	// 循环进行单位转换，直到大小小于 1024
	for floatSize >= 1024 && unitIndex < maxUnixIndex {
		floatSize /= 1024
		unitIndex++
	}
	// 根据不同情况格式化输出
	if unitIndex == 0 {
		return fmt.Sprintf("%d%s", int(floatSize), unitMap[unitIndex])
	}
	return fmt.Sprintf("%.3f%s", floatSize, unitMap[unitIndex])
}
