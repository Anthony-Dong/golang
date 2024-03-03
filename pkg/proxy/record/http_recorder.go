package record

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/anthony-dong/golang/pkg/utils"

	"github.com/valyala/fasthttp"

	"github.com/anthony-dong/golang/pkg/codec/http_codec"
	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/proxy/model"
)

type HttpRecorder struct {
	Storage Storage
	ID      uint64
}

func NewHTTPRecorder(storage Storage) *HttpRecorder {
	return &HttpRecorder{Storage: storage}
}

type ExtraInfo struct {
	IsTLS bool
	Start time.Time
	Spend time.Duration
	Src   net.Addr
}

func (h *HttpRecorder) Record(ctx context.Context, request *fasthttp.Request, response *fasthttp.Response, extra *ExtraInfo) error {
	id := strconv.FormatUint(atomic.AddUint64(&h.ID, 1), 10)
	scheme := "http"
	if extra.IsTLS {
		scheme = "https"
	}
	if h.Storage == nil || isLocalStorage(h.Storage) {
		logs.CtxInfo(ctx, "[ID-%s] [%s] [%s://%s%s] [%d] [%s]", id, request.Header.Method(), scheme, request.Header.Host(), request.URI().Path(), response.StatusCode(), extra.Spend)
		if h.Storage == nil {
			return nil
		}
	}
	req := model.Request{
		Scheme:   scheme,
		Protocol: string(request.Header.Protocol()),
		Method:   string(request.Header.Method()),
		Host:     string(request.Header.Host()),
		Path:     string(request.URI().Path()),
	}
	request.Header.VisitAll(func(key, value []byte) {
		req.Header = append(req.Header, &model.KV{Key: string(key), Value: string(value)})
	})
	request.URI().QueryArgs().VisitAll(func(key, value []byte) {
		req.Query = append(req.Query, &model.KV{Key: string(key), Value: string(value)})
	})
	if cl := request.Header.ContentLength(); cl > 0 && cl <= 1024*128 {
		req.RawBody = request.Body()
	}
	if len(req.RawBody) > 0 {
		body, _ := http_codec.DecodeHttpBody(bytes.NewBuffer(req.RawBody), req.Header, true)
		req.Body = string(body)
	}
	resp := model.Response{
		Protocol:   string(response.Header.Protocol()),
		StatusCode: int32(response.StatusCode()),
	}
	if cl := response.Header.ContentLength(); cl > 0 && cl <= 1024*128 {
		resp.RawBody = response.Body()
	}
	response.Header.VisitAll(func(key, value []byte) {
		resp.Header = append(resp.Header, &model.KV{Key: string(key), Value: string(value)})
	})
	if len(resp.RawBody) > 0 {
		body, _ := http_codec.DecodeHttpBody(bytes.NewBuffer(resp.RawBody), resp.Header, true)
		resp.Body = string(body)
	}
	hp := model.HTTPCapture{
		ID:      id,
		Time:    extra.Start,
		Spend:   int64(extra.Spend / time.Microsecond),
		SrcAddr: extra.Src.String(),
		IsTLS:   extra.IsTLS,

		Request:  &req,
		Response: &resp,
	}
	if isConsulStorage(h.Storage) {
		return h.Storage.Write(utils.String2Bytes(consulLogHTTPCapture(&hp)))
	}
	return h.Storage.Write(utils.ToJsonByte(hp, true))
}

func consulLogHTTPCapture(c *model.HTTPCapture) string {
	scheme := "http"
	if c.IsTLS {
		scheme = "https"
	}
	builder := strings.Builder{}
	builder.WriteString(consulLogSummary("[ID-%s] [%s] [%s] [%s://%s%s] [%d] [%s]\n", c.ID, c.Time.Format("15:04:05.000"), c.Request.Method, scheme, c.Request.Host, c.Request.Path, c.Response.StatusCode, time.Duration(time.Microsecond*time.Duration(c.Spend))))
	consulLogRequest(c.Request, &builder)
	consulLogResponse(c.Response, &builder)
	return builder.String()
}

func consulLogResponse(r *model.Response, builder *strings.Builder) string {
	builder.WriteString(fmt.Sprintf("< %s %d\n", r.Protocol, r.StatusCode))
	for _, elem := range r.Header {
		builder.WriteString(fmt.Sprintf("< %s: %s\n", elem.Key, elem.Value))
	}
	builder.WriteString("< \n")
	if len(r.Body) > 0 {
		builder.WriteString(r.Body)
		builder.WriteString("\n")
	}
	return builder.String()
}

func consulLogRequest(r *model.Request, builder *strings.Builder) {
	values := url.Values{}
	for _, elem := range r.Query {
		values.Add(elem.Key, elem.Value)
	}
	path := r.Path
	if encode := values.Encode(); encode != "" {
		path = path + "?" + encode
	}
	builder.WriteString(fmt.Sprintf("> %s %s %s\n", r.Method, path, r.Protocol))
	for _, elem := range r.Header {
		builder.WriteString(fmt.Sprintf("> %s: %s\n", elem.Key, elem.Value))
	}
	builder.WriteString("> \n")
	if len(r.Body) > 0 {
		builder.WriteString(r.Body)
		builder.WriteString("\n")
	}
}
