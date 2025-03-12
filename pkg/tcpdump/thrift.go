package tcpdump

import (
	"context"
	"fmt"

	"github.com/anthony-dong/golang/pkg/utils"

	"github.com/anthony-dong/golang/pkg/codec/thrift_codec"
)

var _ Message = (*ThriftMessage)(nil)

type ThriftMessage struct {
	Msg *thrift_codec.ThriftMessage
}

func NewThriftMessage(msg *thrift_codec.ThriftMessage) *ThriftMessage {
	return &ThriftMessage{
		Msg: msg,
	}
}

func (*ThriftMessage) Type() MessageType {
	return MessageType_Thrift
}

func (t *ThriftMessage) String() string {
	return utils.ToJson(t.Msg, true)
}

var _ Decoder = (*ThriftDecoder)(nil)

type ThriftDecoder struct{}

func NewThriftDecoder() *ThriftDecoder {
	return &ThriftDecoder{}
}

func (d *ThriftDecoder) Decode(ctx context.Context, reader Reader, packet *TcpPacket) (Message, error) {
	protocol, metaInfo, err := thrift_codec.GetProtocol(ctx, reader)
	if err != nil {
		return nil, fmt.Errorf(`decode thrift protocol find err: %s`, err.Error())
	}
	msg, err := thrift_codec.DecodeMessage(ctx, thrift_codec.NewTProtocol(reader, protocol))
	if err != nil {
		return nil, fmt.Errorf(`decode thrift [protocol=%s] message find err: %s`, protocol, err.Error())
	}
	msg.MetaInfo = metaInfo
	msg.Protocol = protocol
	msg.Transport = &thrift_codec.TransportInfo{FromAddr: packet.Src, ToAddr: packet.Dst}
	return NewThriftMessage(msg), nil
}

func (*ThriftDecoder) Name() string {
	return "thrift"
}
