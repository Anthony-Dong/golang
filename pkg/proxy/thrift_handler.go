package proxy

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/anthony-dong/golang/pkg/proxy/record"

	"github.com/anthony-dong/golang/pkg/logs"

	"github.com/anthony-dong/golang/pkg/codec/thrift_codec"
)

type thriftHandler struct {
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	KeepaliveTime time.Duration

	DialAddr Addr
	Record   *record.ThriftRecorder
}

func NewThriftHandler(dial string, storage record.Storage) *thriftHandler {
	return &thriftHandler{
		KeepaliveTime: time.Second * 15,
		ReadTimeout:   time.Minute * 3,
		WriteTimeout:  time.Minute * 3,
		Record:        record.NewThriftRecorder(storage),
		DialAddr:      MustParseAddr(dial),
	}
}

func (t *thriftHandler) HandlerConn(readConn net.Conn) error {
	writeConn, err := t.DialAddr.Dial()
	if err != nil {
		return err
	}
	defer writeConn.Close()

	ctx := context.Background()
	updateReadTimeout := func(conn net.Conn) {
		if t.ReadTimeout > 0 {
			_ = conn.SetReadDeadline(time.Now().Add(t.ReadTimeout))
		}
	}
	updateWriteTimeout := func(conn net.Conn) {
		if t.ReadTimeout > 0 {
			_ = conn.SetWriteDeadline(time.Now().Add(t.ReadTimeout))
		}
	}
	count := 0
	reader := bufio.NewReader(readConn)
	writer := bufio.NewReader(writeConn)
	for {
		if count > 0 {
			if t.KeepaliveTime <= 0 {
				return nil
			}
			_ = readConn.SetReadDeadline(time.Now().Add(t.KeepaliveTime))
			if _, err := reader.Peek(1); err != nil {
				if err == io.EOF {
					return nil
				}
				logs.CtxWarn(ctx, "conn [%s] find err: %v", readConn.RemoteAddr(), err)
				return err
			}
		}

		start := time.Now()

		updateReadTimeout(readConn)
		req, reqBuffer, err := t.ReadMessage(reader)
		if err != nil {
			return err
		}

		updateWriteTimeout(writeConn)
		if _, err := writeConn.Write(reqBuffer); err != nil {
			return err
		}

		updateReadTimeout(writeConn)
		resp, respBuffer, err := t.ReadMessage(writer)
		if err != nil {
			return err
		}

		updateWriteTimeout(readConn)
		if _, err := readConn.Write(respBuffer); err != nil {
			return err
		}

		_ = t.Record.Record(ctx, req, resp, &record.ThriftExtraInfo{
			SrcAddr: readConn.RemoteAddr(),
			DstAddr: writeConn.RemoteAddr(),
			Time:    start,
			Spend:   time.Now().Sub(start),
		})
		count = count + 1
	}
}

func (t *thriftHandler) ReadMessage(_reader *bufio.Reader) (*thrift_codec.ThriftMessage, []byte, error) {
	reader := &copyReader{reader: _reader}
	ctx := thrift_codec.InjectMateInfo(context.Background())
	protocol, err := thrift_codec.GetProtocol(ctx, reader)
	if err != nil {
		return nil, nil, fmt.Errorf("decode thrift protocol error: %v", err)
	}
	result, err := thrift_codec.DecodeMessage(ctx, thrift_codec.NewTProtocol(reader, protocol))
	if err != nil {
		return nil, nil, fmt.Errorf("decode thrift message error: %v, protocol: %s", err, protocol)
	}
	result.MetaInfo = thrift_codec.GetMateInfo(ctx)
	result.Protocol = protocol
	return result, reader.buffer.Bytes(), nil
}
