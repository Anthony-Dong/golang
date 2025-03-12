package tcpdump

import (
	"fmt"
	"strings"

	"github.com/anthony-dong/golang/pkg/utils"
)

type MessageWriter interface {
	Write(msg Message)
}

var _ MessageWriter = (*consoleLogMessageWriter)(nil)

type consoleLogMessageWriter struct {
	enable map[MessageType]bool
}

func NewConsoleLogMessageWriter(enable []MessageType) MessageWriter {
	enableMap := make(map[MessageType]bool)
	for _, m := range enable {
		enableMap[m] = true
	}
	return &consoleLogMessageWriter{enable: enableMap}

}

func (c *consoleLogMessageWriter) Write(msg Message) {
	if c.enable[msg.Type()] {
		ss := utils.ToString(msg)
		fmt.Println(strings.TrimSpace(ss))
	}
}
