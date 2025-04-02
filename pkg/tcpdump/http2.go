package tcpdump

import (
	"bufio"
	"context"
	"net/http"

	http2_codec "github.com/anthony-dong/golang/pkg/codec/http2"
	"github.com/anthony-dong/golang/pkg/utils"
	"golang.org/x/net/http2"
)

var _ Decoder = (*Http2Decoder)(nil)

type Http2Decoder struct {
}

func NewHttp2Decoder() Decoder {
	return &Http2Decoder{}
}

func (h *Http2Decoder) Name() string {
	return "http2"
}

func (h *Http2Decoder) Decode(ctx context.Context, reader Reader, packet *TcpPacket) (Message, error) {
	r := bufio.NewReader(reader)
	if message, err := http2_codec.DecodeH2CMessage(r); err == nil {
		return &H2cMessage{Req: message}, nil
	}
	result := make([]Message, 0)
	for {
		frame, err := http2_codec.DecodeFrame(r)
		if err != nil {
			if len(result) == 0 {
				return nil, err
			}
			break
		}
		http2Frame := NewHttp2Frame(frame)
		http2Frame.Transport = &http2_codec.TransportInfo{
			FromAddr: packet.Src,
			ToAddr:   packet.Dst,
		}
		result = append(result, http2Frame)
	}
	return NewMultiMessage(result), nil
}

type H2cMessage struct {
	Req *http.Request
}

func (*H2cMessage) Type() MessageType {
	return MessageType_HTTP2
}

func (*H2cMessage) String() string {
	frame := http2_codec.HTTP2Frame{
		Type: "h2c",
		Data: "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n",
	}
	return utils.ToJson(frame, true)
}

type Http2Frame struct {
	Frame     http2.Frame
	Transport *http2_codec.TransportInfo
}

func (*Http2Frame) Type() MessageType {
	return MessageType_HTTP2
}

func (f *Http2Frame) String() string {
	frame := http2_codec.ConvertToFrame(f.Frame)
	frame.TransportInfo = f.Transport
	return utils.ToJson(frame, true)
}

func NewHttp2Frame(frame http2.Frame) *Http2Frame {
	return &Http2Frame{Frame: frame}
}
