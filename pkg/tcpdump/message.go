package tcpdump

import (
	"fmt"
	"time"

	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/utils"
)

type Message interface {
	Type() MessageType
}

type MessageType string

const (
	MessageType_Unknown   = "Unknown"
	MessageType_Log       = "Log"
	MessageType_TcpPacket = "TcpPacket"
	MessageType_HTTP      = "HTTP"
	MessageType_Thrift    = "Thrift"
	MessageType_Layer     = "Layer"

	MessageType_TcpdumpHeader  = "TcpdumpHeader"
	MessageType_TcpdumpPayload = "TcpdumpPayload"
)

var _ Message = (*UnknownMessage)(nil)
var _ Message = (*LogMessage)(nil)

type LogMessage struct {
	Time  time.Time  `json:"time"`
	Level logs.Level `json:"level"`
	Msg   string     `json:"msg"`
}

func NewLogMessage(level logs.Level, format string, v ...interface{}) *LogMessage {
	return &LogMessage{
		Time:  time.Now(),
		Level: level,
		Msg:   fmt.Sprintf(format, v...),
	}
}

func (*LogMessage) Type() MessageType {
	return MessageType_Log
}

func (m *LogMessage) MarshalJSON() ([]byte, error) {
	return utils.ToJsonByte(m), nil
}

type UnknownMessage struct {
}

func (*UnknownMessage) Type() MessageType {
	return MessageType_Unknown
}

type TcpdumpHeader struct {
	Header string
}

func (*TcpdumpHeader) Type() MessageType {
	return MessageType_TcpdumpHeader
}

func (m *TcpdumpHeader) String() string {
	return m.Header
}

type TcpdumpPayload struct {
	Payload []byte
}

func (*TcpdumpPayload) Type() MessageType {
	return MessageType_TcpdumpPayload
}

func (m *TcpdumpPayload) String() string {
	return HexDump(m.Payload)
}
