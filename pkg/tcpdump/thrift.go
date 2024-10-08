package tcpdump

import (
	"fmt"

	"github.com/anthony-dong/golang/pkg/utils"

	"github.com/pkg/errors"

	"github.com/anthony-dong/golang/pkg/codec/thrift_codec"
)

func NewThriftDecoder() Decoder {
	return func(ctx *Context, reader SourceReader) error {
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
		ctx.PrintPayload(utils.ToJson(result, true))
		return nil
	}
}
