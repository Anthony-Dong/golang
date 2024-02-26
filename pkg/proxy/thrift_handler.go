package proxy

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/utils"

	"github.com/anthony-dong/golang/pkg/codec/thrift_codec"
)

type thriftHandler struct {
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	KeepaliveTime time.Duration

	Recorder func(payload interface{}, isReq bool)
}

func NewThriftHandler(r func(payload interface{}, isReq bool)) *thriftHandler {
	return &thriftHandler{
		KeepaliveTime: time.Second * 15,
		ReadTimeout:   time.Minute * 3,
		WriteTimeout:  time.Minute * 3,
		Recorder:      r,
	}
}

func (t *thriftHandler) HandlerConn(readConn net.Conn, dialAddr Addr) error {
	writeConn, err := dialAddr.Dial()
	if err != nil {
		return err
	}
	ctx := context.Background()
	defer func() {
		if err := writeConn.Close(); err != nil {
			logs.CtxWarn(ctx, "conn [%s] close find err: %v", writeConn.RemoteAddr(), err)
			return
		}
	}()
	if t.Recorder == nil {
		t.Recorder = nonRecorder
	}
	count := 0
	reader := bufio.NewReader(readConn)
	writer := bufio.NewReader(writeConn)
	for {
		if count > 0 {
			// keep alive 10s
			if t.KeepaliveTime > 0 {
				_ = readConn.SetReadDeadline(time.Now().Add(t.KeepaliveTime))
			}
			if _, err := reader.Peek(1); err != nil {
				if err == io.EOF {
					logs.CtxDebug(ctx, "conn [%s] close", readConn.RemoteAddr())
					return nil
				}
				logs.CtxWarn(ctx, "conn [%s] find err: %v", readConn.RemoteAddr(), err)
				return err
			}
		}
		if t.ReadTimeout > 0 {
			_ = readConn.SetReadDeadline(time.Now().Add(t.ReadTimeout))
		}
		req, reqBuffer, err := t.ReadMessage(reader)
		if err != nil {
			return err
		}
		t.Recorder(req, true)
		if t.WriteTimeout > 0 {
			_ = writeConn.SetWriteDeadline(time.Now().Add(t.WriteTimeout))
		}
		if _, err := writeConn.Write(reqBuffer); err != nil {
			return err
		}
		if t.ReadTimeout > 0 {
			_ = writeConn.SetReadDeadline(time.Now().Add(t.ReadTimeout))
		}
		resp, respBuffer, err := t.ReadMessage(writer)
		if err != nil {
			return err
		}
		t.Recorder(resp, false)
		if t.WriteTimeout > 0 {
			_ = readConn.SetWriteDeadline(time.Now().Add(t.WriteTimeout))
		}
		if _, err := readConn.Write(respBuffer); err != nil {
			return err
		}
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

func nonRecorder(payload interface{}, isReq bool) {
}

var consoleRecorderLock sync.Mutex

func ConsoleRecorder(payload interface{}, isReq bool) {
	consoleRecorderLock.Lock()
	defer consoleRecorderLock.Unlock()
	if isReq {
		logs.StdOut(strings.Repeat(">", 30) + " [CALL] " + time.Now().Format(utils.FormatTimeV1) + " " + strings.Repeat(">", 30))
	} else {
		logs.StdOut(strings.Repeat("<", 29) + " [REPLY] " + time.Now().Format(utils.FormatTimeV1) + " " + strings.Repeat("<", 30))
	}
	logs.StdOut(utils.ToJson(payload, true))
}
