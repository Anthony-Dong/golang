package record

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/fatih/color"

	"github.com/anthony-dong/golang/pkg/codec/thrift_codec"
	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/utils"
)

type ThriftRecorder struct {
	ID      uint64
	Storage Storage
}

func NewThriftRecorder(storage Storage) *ThriftRecorder {
	return &ThriftRecorder{Storage: storage}
}

type ThriftExtraInfo struct {
	ID      string
	SrcAddr net.Addr
	DstAddr net.Addr
	Time    time.Time
	Spend   time.Duration
}

func (t *ThriftRecorder) Record(ctx context.Context, req *thrift_codec.ThriftMessage, resp *thrift_codec.ThriftMessage, extra *ThriftExtraInfo) error {
	id := atomic.AddUint64(&t.ID, 1)
	extra.ID = strconv.Itoa(int(id))
	if t.Storage == nil || isLocalStorage(t.Storage) {
		status := "OK"
		if resp.MessageType == thrift_codec.EXCEPTION {
			status = "EXCEPTION"
		}
		logs.CtxInfo(ctx, "[ID-%d] %s->%s [%s] [%s] [%s]", extra.ID, extra.SrcAddr, extra.DstAddr, req.Method, status, extra.Spend)
		if t.Storage == nil {
			return nil
		}
	}
	if isConsulStorage(t.Storage) {
		return t.Storage.Write(utils.String2Bytes(consulLogThriftMessage(req, resp, extra)))
	}
	return t.Storage.Write(utils.ToJsonByte(map[string]interface{}{
		"request":  req,
		"response": resp,
		"src_addr": extra.SrcAddr.String(),
		"dst_addr": extra.DstAddr.String(),
		"time":     extra.Time,
		"spend":    extra.Spend / time.Microsecond,
	}, true))
}

func consulLogSummary(format string, v ...interface{}) string {
	return color.HiGreenString(format, v...)
}

func consulLogThriftMessage(req *thrift_codec.ThriftMessage, resp *thrift_codec.ThriftMessage, extra *ThriftExtraInfo) string {
	toJson := func(v interface{}) string {
		indent, _ := json.MarshalIndent(v, "", "\t")
		return string(indent)
	}

	status := "OK"
	if resp.MessageType == thrift_codec.EXCEPTION {
		status = "EXCEPTION"
	}

	builder := strings.Builder{}
	builder.WriteString(consulLogSummary("[ID-%s] [%s] %s->%s [%s] [%s] [%s]\n", extra.ID, extra.Time.Format("15:04:05.000"), extra.SrcAddr, extra.DstAddr, req.Method, status, extra.Spend))

	builder.WriteString(fmt.Sprintf("> %s %s %s\n", strings.ToUpper(req.MessageType.String()), req.Method, req.Protocol))
	builder.WriteString(fmt.Sprintf("> ID: %d\n", req.SeqId))
	if req.MetaInfo != nil {
		for _, elem := range utils.SortMap(req.MetaInfo.StrInfo, func(i, j string) bool { return i < j }) {
			builder.WriteString(fmt.Sprintf("> %s: %s\n", elem.Key, elem.Value))
		}
		for _, elem := range utils.SortMap(req.MetaInfo.IntInfo, func(i, j uint16) bool { return i < j }) {
			builder.WriteString(fmt.Sprintf("> @%d: %s\n", elem.Key, elem.Value))
		}
	}
	builder.WriteString(">\n")
	builder.WriteString(toJson(req.Payload))
	builder.WriteString("\n")

	builder.WriteString(fmt.Sprintf("< %s %s %s\n", strings.ToUpper(resp.MessageType.String()), resp.Method, resp.Protocol))
	builder.WriteString(fmt.Sprintf("< ID: %d\n", req.SeqId))
	if resp.MetaInfo != nil {
		for _, elem := range utils.SortMap(resp.MetaInfo.StrInfo, func(i, j string) bool { return i < j }) {
			builder.WriteString(fmt.Sprintf("> %s: %s\n", elem.Key, elem.Value))
		}
		for _, elem := range utils.SortMap(resp.MetaInfo.IntInfo, func(i, j uint16) bool { return i < j }) {
			builder.WriteString(fmt.Sprintf("> @%d: %s\n", elem.Key, elem.Value))
		}
	}
	builder.WriteString("<\n")
	switch resp.MessageType {
	case thrift_codec.EXCEPTION:
		builder.WriteString(toJson(resp.Exception))
	case thrift_codec.REPLY:
		builder.WriteString(toJson(resp.Payload))
	}
	builder.WriteString("\n")

	return builder.String()
}
