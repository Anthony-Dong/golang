package proxy

import (
	"time"

	"github.com/valyala/fasthttp"

	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/proxy/record"
)

func NewRecordHTTPHandler(storage record.Storage) fasthttp.RequestHandler {
	httpConfig := NewDefaultHttpConfig()
	client := NewHTTPClient(httpConfig)
	recorder := record.NewHTTPRecorder(storage)
	return func(ctx *fasthttp.RequestCtx) {
		scheme := "http"
		isTls, _ := ctx.Value("is_tls").(bool)
		if isTls {
			scheme = "https"
		}
		request := &ctx.Request
		response := &ctx.Response
		start := time.Now()
		if err := client.Do(request, response); err != nil {
			logs.CtxError(ctx, "[%s] [%s] [%s] %s err: %v", scheme, request.Header.Method(), request.Header.Host(), request.Header.RequestURI(), err)
			returnError(ctx, err, `Proxy Request Error`)
			return
		}
		_ = recorder.Record(ctx, request, response, &record.ExtraInfo{IsTLS: isTls, Start: start, Spend: time.Now().Sub(start), Src: ctx.RemoteAddr()})
	}
}

type HTTPBizError struct {
	HTTPCode int `json:"http_code"`
}
