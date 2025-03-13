package tcpdump

import (
	"fmt"
	"strings"
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
		fmt.Println(strings.TrimSpace(msg.String()))
	}
}
