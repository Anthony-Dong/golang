package thrift_codec

import (
	"context"
	"encoding/binary"
	"errors"
	"io"

	"github.com/anthony-dong/golang/pkg/utils"

	"github.com/apache/thrift/lib/go/thrift"

	"github.com/anthony-dong/golang/pkg/codec/thrift_codec/kitex"
)

type Protocol uint8

// Unframed 又成为 Buffered协议.
const (
	UnknownProto Protocol = iota

	// Unframed协议大改分为以下几类.
	UnframedBinary
	UnframedCompact

	// Framed协议分为以下几类.
	FramedBinary
	FramedCompact

	// Header 协议，默认是Unframed，也可以是Framed的，其实本身来说Header协议并不需要再包一层Framed协议.
	UnframedHeader
	FramedHeader

	// Binary 非严格协议大改分为以下两种！其实还有一种是 Header+Binary这种协议，这里就不做细份了.
	UnframedUnStrictBinary
	FramedUnStrictBinary

	// kitex protocol
	UnframedBinaryTTHeader
	FramedBinaryTTHeader

	UnframedBinaryMeshHeader
	FramedBinaryMeshHeader
)

type reader interface {
	io.Reader
	Peek(int) ([]byte, error)
}

func NewTProtocol(reader io.Reader, protocol Protocol) thrift.TProtocol {
	tReader := thrift.NewStreamTransportR(reader)
	return newTProtocol(tReader, protocol, false)
}

func newTProtocol(stream *thrift.StreamTransport, protocol Protocol, isW bool) thrift.TProtocol {
	switch protocol {
	case UnframedBinary, UnframedBinaryTTHeader, UnframedBinaryMeshHeader:
		return thrift.NewTBinaryProtocolTransport(stream)
	case UnframedUnStrictBinary:
		strictWrite := true
		if isW {
			strictWrite = false
		}
		return thrift.NewTBinaryProtocol(stream, false, strictWrite)
	case UnframedCompact:
		return thrift.NewTCompactProtocol(stream)
	case FramedBinary, FramedBinaryTTHeader, FramedBinaryMeshHeader:
		return thrift.NewTBinaryProtocolTransport(thrift.NewTFramedTransport(stream))
	case FramedUnStrictBinary:
		strictWrite := true
		if isW {
			strictWrite = false
		}
		return thrift.NewTBinaryProtocol(thrift.NewTFramedTransport(stream), false, strictWrite)
	case FramedCompact:
		return thrift.NewTCompactProtocol(thrift.NewTFramedTransport(stream))
	case UnframedHeader:
		return thrift.NewTHeaderProtocol(stream)
	case FramedHeader:
		return thrift.NewTHeaderProtocol(thrift.NewTFramedTransport(stream))
	default:
		return thrift.NewTBinaryProtocolTransport(stream)
	}
}

func NewTProtocolEncoder(writer io.Writer, protocol Protocol) thrift.TProtocol {
	tWriter := thrift.NewStreamTransportW(writer)
	return newTProtocol(tWriter, protocol, true)
}

func readBytes(r io.Reader, len int) ([]byte, error) {
	result := make([]byte, len)
	if _, err := r.Read(result); err != nil {
		return nil, err
	}
	return result, nil
}

// flag 为4字节
func IsUnframedBinary(reader reader, offset int) bool {
	/**
	Binary protocol Message, strict encoding, 12+ bytes:
	+--------+--------+--------+--------+--------+--------+--------+--------+--------+...+--------+--------+--------+--------+--------+
	|1vvvvvvv|vvvvvvvv|unused  |00000mmm| name length                       | name                | seq id                            |
	+--------+--------+--------+--------+--------+--------+--------+--------+--------+...+--------+--------+--------+--------+--------+
	*/
	flag, err := reader.Peek(offset + Size32)
	if err != nil {
		return false
	}
	flag = flag[offset:]
	// 取前两个字节版本号，如果 VERSION_1 = 0x80010000
	return binary.BigEndian.Uint32(flag)&thrift.VERSION_MASK == thrift.VERSION_1
}

const (
	Size32 = 4
	Size16 = 2
	Size8  = 1

	FrameHeaderSize = 4
)

func IsUnframedUnStrictBinary(reader reader, offset int) bool {
	/**
	UnframedBinary 的非严格模式，头部4字节一定会大于0
	Binary protocol Message, old encoding, 9+ bytes:
	+--------+--------+--------+--------+--------+...+--------+--------+--------+--------+--------+--------+
	| name length                       | name                |00000mmm| seq id                            |
	+--------+--------+--------+--------+--------+...+--------+--------+--------+--------+--------+--------+
	name length(四字节): 这里为了兼容上面协议一，所以高位第一个bit必须为0！也就是name length必须要有符号的正数！
	*/
	flag, err := reader.Peek(offset + Size32)
	if err != nil {
		return false
	}
	flag = flag[offset:]

	if nameLen := binary.BigEndian.Uint32(flag); nameLen > 0 {
		headSize := int(Size32 + nameLen + Size8)
		headBuf, err := reader.Peek(offset + headSize)
		if err != nil {
			return false
		}
		headBuf = headBuf[offset:]

		return headBuf[headSize-1]&0xf8 == 0 && thrift.TMessageType(headBuf[headSize-1]) <= thrift.ONEWAY
	}
	return false
}

func IsUnframedCompact(reader reader, offset int) bool {
	/**
	Compact protocol Message (4+ bytes):
	+--------+--------+--------+...+--------+--------+...+--------+--------+...+--------+
	|pppppppp|mmmvvvvv| seq id              | name length         | name                |
	+--------+--------+--------+...+--------+--------+...+--------+--------+...+--------+
	*/
	flag, err := reader.Peek(offset + Size32)
	if err != nil {
		return false
	}
	flag = flag[offset:]
	return flag[0] == thrift.COMPACT_PROTOCOL_ID && flag[1]&thrift.COMPACT_VERSION_MASK == thrift.COMPACT_VERSION
}

func IsFramedBinary(reader reader, offset int) bool {
	return IsUnframedBinary(reader, offset+4)
}

func IsUnframedHeader(reader reader, offset int) bool {
	/**
	THeader proto
	  0 1 2 3 4 5 6 7 8 9 a b c d e f 0 1 2 3 4 5 6 7 8 9 a b c d e f
	+----------------------------------------------------------------+
	| 0|                          LENGTH                             |
	+----------------------------------------------------------------+
	| 0|       HEADER MAGIC          |            FLAGS              |
	+----------------------------------------------------------------+
	*/
	flag, err := reader.Peek(offset + Size32*2)
	if err != nil {
		return false
	}
	flag = flag[offset:]
	if binary.BigEndian.Uint32(flag[Size32:])&thrift.THeaderHeaderMask == thrift.THeaderHeaderMagic {
		if binary.BigEndian.Uint32(flag[:Size32]) > thrift.THeaderMaxFrameSize {
			return false
			//return UnknownProto, thrift.NewTProtocolExceptionWithType(
			//	thrift.SIZE_LIMIT,
			//	errors.New("frame too large"),
			//)
		}
		return true
	}
	return false
}

// GetProtocol 自动获取请求的消息协议！记住一定是消息协议！
// reader *bufio.Reader 类型是因为重复读！
// GetProtocol 使用前需要通过 InjectMateInfo 注入MetaInfo
func GetProtocol(ctx context.Context, reader reader) (Protocol, *kitex.MetaInfo, error) {
	if IsUnframedHeader(reader, 0) {
		return UnframedHeader, nil, nil
	}
	if IsUnframedHeader(reader, FrameHeaderSize) {
		return FramedHeader, nil, nil
	}
	if IsUnframedBinary(reader, 0) {
		return UnframedBinary, nil, nil
	}
	if IsUnframedBinary(reader, FrameHeaderSize) {
		return FramedBinary, nil, nil
	}
	if IsUnframedCompact(reader, 0) {
		return UnframedCompact, nil, nil
	}
	if IsUnframedCompact(reader, FrameHeaderSize) {
		return FramedCompact, nil, nil
	}
	if kitex.IsTTHeader(reader) {
		metaInfo := &kitex.MetaInfo{}
		size, err := kitex.ReadTTHeader(reader, metaInfo)
		if err != nil {
			return UnknownProto, nil, err
		}
		if IsUnframedBinary(reader, size) {
			_ = utils.SkipReader(reader, size)
			return UnframedBinaryTTHeader, metaInfo, nil
		}
		if IsFramedBinary(reader, size) {
			_ = utils.SkipReader(reader, size)
			return FramedBinaryTTHeader, metaInfo, nil
		}
	}
	if kitex.IsMeshHeader(reader) {
		metaInfo := &kitex.MetaInfo{}
		size, err := kitex.ReadMeshHeader(reader, metaInfo)
		if err != nil {
			return UnknownProto, nil, err
		}
		if IsUnframedBinary(reader, size) {
			if err := utils.SkipReader(reader, size); err != nil {
				return UnknownProto, nil, err
			}
			return UnframedBinaryMeshHeader, metaInfo, nil
		}
		if IsFramedBinary(reader, size) {
			if err := utils.SkipReader(reader, size); err != nil {
				return UnknownProto, nil, err
			}
			return FramedBinaryMeshHeader, metaInfo, nil
		}
	}
	// UnStrict protocol 这个需要读取的buffer比较大(遇到buffer不满足的时候会长时间block)，目前没有太好的优化手动，因为放在最后面了
	// 还有一种解决方案就是限制name_len的大小，理论上name不会有很长，撑死了也就100个字符左右，但是遇到小包也会阻塞
	if IsUnframedUnStrictBinary(reader, 0) {
		return UnframedUnStrictBinary, nil, nil
	}
	if IsUnframedUnStrictBinary(reader, FrameHeaderSize) {
		return FramedUnStrictBinary, nil, nil
	}
	return UnknownProto, nil, thrift.NewTProtocolExceptionWithType(
		thrift.UNKNOWN_PROTOCOL_EXCEPTION,
		errors.New("unknown protocol"),
	)
}
