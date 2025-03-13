package tcpdump

import (
	"fmt"
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

func ConvertToMessageType(input []string) []MessageType {
	result := make([]MessageType, len(input))
	for i, s := range input {
		result[i] = MessageType(s)
	}
	return result
}

func (c *consoleLogMessageWriter) Write(msg Message) {
	if c.enable[msg.Type()] {
		output := msg.String()
		if len(output) > 0 && output[len(output)-1] == '\n' {
			fmt.Printf(output)
		} else {
			fmt.Println(output)
		}
	}
}
