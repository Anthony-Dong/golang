package tcpdump

import (
	"context"

	"github.com/anthony-dong/golang/pkg/logs"
)

type Message interface {
	Type() MessageType
	String() string
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
	Msg string `json:"msg"`
}

func NewLogMessage(ctx context.Context, level logs.Level, format string, v ...interface{}) *LogMessage {
	return &LogMessage{
		Msg: logs.Sprintf(ctx, logs.GetFlag(), level, -1, format, v...),
	}
}

func (*LogMessage) Type() MessageType {
	return MessageType_Log
}

func (m *LogMessage) String() string {
	return m.Msg
}

type UnknownMessage struct {
}

func (*UnknownMessage) Type() MessageType {
	return MessageType_Unknown
}

func (*UnknownMessage) String() string {
	return "<UnknownMessage>"
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
