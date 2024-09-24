package rpc

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"

	"github.com/anthony-dong/golang/pkg/codec/thrift_codec"
	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
)

type thriftServer struct {
	addr string
	desc *descriptor.ServiceDescriptor

	listen net.Listener
}

func NewThriftServer(desc *descriptor.ServiceDescriptor, addr string) *thriftServer {
	return &thriftServer{desc: desc, addr: addr}
}

func (t *thriftServer) encode(oprot thrift.TProtocol, method string, seq int32) error {
	function := t.desc.Functions[method]
	if function == nil {
		return fmt.Errorf("function %s not found", method)
	}
	value, err := GetThriftExampleValue(function.Response, nil, nil)
	if err != nil {
		return err
	}
	if err := thrift_codec.EncodeReply(oprot, function, seq, value); err != nil {
		return err
	}
	return nil
}

func (t *thriftServer) Run() error {
	listen, err := net.Listen("tcp", t.addr)
	if err != nil {
		return err
	}
	t.listen = listen
	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}
		go func() {
			ctx := context.Background()
			if err := t.handle(ctx, conn); err != nil {
				logs.CtxError(ctx, "handle conn %s find err: %v", conn.RemoteAddr(), err)
			}
			_ = conn.Close()
		}()
	}
}

func (t *thriftServer) Close() error {
	if t.listen != nil {
		return t.listen.Close()
	}
	return nil
}

func (t *thriftServer) handle(ctx context.Context, conn net.Conn) error {
	readBuffer := bufio.NewReader(conn)
	for {
		protocol, metaInfo, err := thrift_codec.GetProtocol(ctx, readBuffer)
		if err != nil {
			if _, err := readBuffer.Peek(1); err != nil {
				if err == io.EOF {
					return nil
				}
			}
			return err
		}
		tProtocol := thrift_codec.NewTProtocol(readBuffer, protocol)
		request, err := thrift_codec.DecodeMessage(ctx, tProtocol)
		if err != nil {
			return err
		}
		request.MetaInfo = metaInfo
		if request.MessageType == thrift_codec.ONEWAY { // oneway skip
			continue
		}
		encoder := thrift_codec.NewTProtocolEncoder(conn, protocol)
		if err := t.encode(encoder, request.Method, request.SeqId); err != nil {
			return err
		}
		if err := encoder.Flush(ctx); err != nil {
			return err
		}
	}
}
