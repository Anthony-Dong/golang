package tcpdump

import (
	"context"
	"fmt"

	"github.com/anthony-dong/golang/pkg/utils"

	"github.com/pkg/errors"

	"github.com/anthony-dong/golang/pkg/codec/thrift_codec"
)

func NewThriftDecoder(parser ThriftMessageParser) Decoder {
	return func(ctx *Context, reader SourceReader, packet Packet) error {
		protocol, metaInfo, err := thrift_codec.GetProtocol(ctx, reader)
		if err != nil {
			return errors.Wrap(err, "decode thrift protocol error")
		}
		result, err := thrift_codec.DecodeMessage(ctx, thrift_codec.NewTProtocol(reader, protocol))
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("decode thrift message error, protocol: %s", protocol))
		}
		result.MetaInfo = metaInfo
		result.Protocol = protocol
		result.Transport = &thrift_codec.TransportInfo{FromAddr: packet.Src, ToAddr: packet.Dst}
		message, err := parser.ParseMessage(ctx, result, packet)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("parse message error, protocol: %s", protocol))
		}
		ctx.PrintPayload(utils.Bytes2String(message))
		return nil
	}
}

type ThriftMessageParser interface {
	ParseMessage(ctx context.Context, msg *thrift_codec.ThriftMessage, _ Packet) ([]byte, error)
}

type defaultThriftMessageParser struct {
}

func NewThriftMessageParser() ThriftMessageParser {
	return &defaultThriftMessageParser{}
}

func (*defaultThriftMessageParser) ParseMessage(ctx context.Context, msg *thrift_codec.ThriftMessage, _ Packet) ([]byte, error) {
	return utils.ToJsonByte(msg, true), nil
}
