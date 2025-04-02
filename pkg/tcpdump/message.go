package tcpdump

import (
	"context"

	"github.com/anthony-dong/golang/pkg/logs"
)

type Message interface {
	Type() MessageType
	String() string
}

type MultiMessage struct {
	Messages []Message
}

func NewMultiMessage(msg []Message) *MultiMessage {
	return &MultiMessage{Messages: msg}
}

func (*MultiMessage) Type() MessageType {
	return MessageType_Multi
}

func (*MultiMessage) String() string {
	return "MultiMessage"
}

type MessageType string

const (
	MessageType_Multi     = "multi"
	MessageType_Unknown   = "unknown"
	MessageType_Log       = "log"
	MessageType_TcpPacket = "tcp"
	MessageType_HTTP      = "http"
	MessageType_Thrift    = "thrift"
	MessageType_Layer     = "layer"
	MessageType_Tcpdump   = "tcpdump"

	MessageType_HTTP2 = "http2"
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
	return MessageType_Tcpdump
}

func (m *TcpdumpHeader) String() string {
	return m.Header
}

type TcpdumpPayload struct {
	Payload []byte
}

func (*TcpdumpPayload) Type() MessageType {
	return MessageType_Tcpdump
}

func (m *TcpdumpPayload) String() string {
	return HexDump(m.Payload)
}
