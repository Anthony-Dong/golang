package tcpdump

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httputil"
	"unsafe"

	"github.com/anthony-dong/golang/pkg/utils"
	"github.com/pkg/errors"

	"github.com/anthony-dong/golang/pkg/codec/http_codec"
)

var _ Message = (*HttpMessage)(nil)

type HttpMessage struct {
	Req  *http.Request
	Resp *http.Response

	Host    string
	RawData []byte
}

func (m *HttpMessage) String() string {
	if m.Req != nil {
		if m.enableDumpRequest(m.Req) {
			if req, err := DumpRequest(m.Req); err == nil {
				return UnsafeString(req)
			}
		}
		return UnsafeString(m.RawData)
	}
	if m.enableDumpResponse(m.Resp) {
		if resp, err := DumpResponse(m.Resp); err == nil {
			return UnsafeString(resp)
		}
	}
	return UnsafeString(m.RawData)
}

func UnsafeString(buf []byte) string {
	return *(*string)(unsafe.Pointer(&buf))
}

func (m *HttpMessage) enableDumpRequest(req *http.Request) bool {
	if encoding := req.Header.Get("Content-Encoding"); encoding != "" {
		return req.ContentLength > 0
	}
	return false
}

func (m *HttpMessage) enableDumpResponse(resp *http.Response) bool {
	if encoding := resp.Header.Get("Content-Encoding"); encoding != "" {
		return resp.ContentLength > 0
	}
	if utils.Contains(resp.TransferEncoding, "chunked") {
		return true
	}
	return false
}

func DumpRequest(req *http.Request) ([]byte, error) {
	body, err := http_codec.DecodeHttpBody(req.Body, req.Header, true)
	if err != nil {
		return nil, err
	}
	dumpRequest, err := httputil.DumpRequest(req, false)
	if err != nil {
		return nil, err
	}
	writer := bytes.NewBuffer(dumpRequest)
	writer.Write(body)
	return writer.Bytes(), nil
}

func DumpResponse(resp *http.Response) ([]byte, error) {
	body, err := http_codec.DecodeHttpBody(resp.Body, resp.Header, true)
	if err != nil {
		return nil, err
	}
	dumpResponse, err := httputil.DumpResponse(resp, false)
	if err != nil {
		return nil, err
	}
	writer := bytes.NewBuffer(dumpResponse)
	writer.Write(body)
	return writer.Bytes(), nil
}

func (*HttpMessage) Type() MessageType {
	return MessageType_HTTP
}

func NewHttpReqMessage(req *http.Request, rawData []byte, host string) *HttpMessage {
	return &HttpMessage{
		Req:     req,
		RawData: rawData,
		Host:    host,
	}
}

func NewHttpRespMessage(resp *http.Response, rawData []byte) *HttpMessage {
	return &HttpMessage{
		Resp:    resp,
		RawData: rawData,
	}
}

var _ Decoder = (*HttpDecoder)(nil)

type HttpDecoder struct{}

func NewHttpDecoder() Decoder {
	return &HttpDecoder{}
}

func (h *HttpDecoder) Decode(ctx context.Context, reader Reader, packet *TcpPacket) (Message, error) {
	crlfNum := 0 // /r/n 换行符， http协议分割符号本质上是换行符！所以清除头部的换行符(假如存在这种case)
	for {
		peek, err := reader.Peek(2 + crlfNum)
		if err != nil {
			return nil, errors.Wrap(err, `read http content error`)
		}
		peek = peek[crlfNum:]
		if peek[0] == '\r' && peek[1] == '\n' {
			crlfNum = crlfNum + 2
			continue
		}
		break
	}
	if crlfNum != 0 {
		if _, err := reader.Read(make([]byte, crlfNum)); err != nil {
			return nil, errors.Wrap(err, `read http content error`)
		}
	}
	copyR := &bytes.Buffer{}
	bufReader := bufio.NewReader(io.TeeReader(reader, copyR)) // copy

	isRequest, err := isHttpRequest(ctx, reader)
	if err != nil {
		return nil, errors.Wrap(err, `read http request content error`)
	}
	if isRequest {
		req, err := http.ReadRequest(bufReader)
		if err != nil {
			return nil, errors.Wrap(err, `read http request content err`)
		}
		return NewHttpReqMessage(req, copyR.Bytes(), packet.Src), nil
	}

	isResponse, err := isHttpResponse(ctx, reader)
	if err != nil {
		return nil, errors.Wrap(err, `read http response content error`)
	}
	if !isResponse {
		return nil, errors.Errorf(`invalid http content`)
	}
	resp, err := http.ReadResponse(bufReader, nil)
	if err != nil {
		return nil, errors.Wrap(err, `read http response content error`)
	}
	if resp.ContentLength > 0 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, `read http response content error`)
		}
		resp.Body = io.NopCloser(bytes.NewBuffer(body))
	} else if utils.Contains(resp.TransferEncoding, "chunked") {
		chunked, err := http_codec.ReadChunked(bufReader)
		if err != nil {
			return nil, errors.Wrap(err, `read http response content error, transfer encoding is chunked`)
		}
		resp.Body = io.NopCloser(bytes.NewBuffer(chunked)) // copy
	}
	return NewHttpRespMessage(resp, copyR.Bytes()), nil
}

func (h *HttpDecoder) Name() string {
	return "http1.1"
}

// 	MethodGet     = "GET"
//	MethodHead    = "HEAD"
//	MethodPost    = "POST"
//	MethodPut     = "PUT"
//	MethodPatch   = "PATCH" // RFC 5789
//	MethodDelete  = "DELETE"
//	MethodConnect = "CONNECT"
//	MethodOptions = "OPTIONS"
//	MethodTrace   = "TRACE"

func isHttpResponse(ctx context.Context, reader Reader) (bool, error) {
	peek, err := reader.Peek(6)
	if err != nil {
		return false, err
	}
	if string(peek) == "HTTP/1" {
		return true, nil
	}
	return false, nil
}
func isHttpRequest(ctx context.Context, reader Reader) (bool, error) {
	peek, err := reader.Peek(7)
	if err != nil {
		return false, err
	}
	if method := string(peek[:3]); method == "GET" || method == "PUT" {
		return true, nil
	}
	if method := string(peek[:4]); method == "HEAD" || method == "POST" {
		return true, nil
	}
	if method := string(peek[:5]); method == "PATCH" || method == "TRACE" {
		return true, nil
	}
	if method := string(peek[:6]); method == "DELETE" {
		return true, nil
	}
	if method := string(peek[:7]); method == "OPTIONS" || method == "CONNECT" {
		return true, nil
	}
	return false, nil
}

//
//func NewHTTP1Decoder() Decoder {
//	return func(ctx *Context, reader SourceReader, _ Packet) error {
//		crlfNum := 0 // /r/n 换行符， http协议分割符号本质上是换行符！所以清除头部的换行符(假如存在这种case)
//		for {
//			peek, err := reader.Peek(2 + crlfNum)
//			if err != nil {
//				return errors.Wrap(err, `read http content error`)
//			}
//			peek = peek[crlfNum:]
//			if peek[0] == '\r' && peek[1] == '\n' {
//				crlfNum = crlfNum + 2
//				continue
//			}
//			break
//		}
//		if crlfNum != 0 {
//			if _, err := reader.Read(make([]byte, crlfNum)); err != nil {
//				return errors.Wrap(err, `read http content error`)
//			}
//		}
//
//		copyR := bufutils.NewBuffer()
//		defer bufutils.ResetBuffer(copyR)
//		bufReader := bufio.NewReader(io.TeeReader(reader, copyR)) // copy
//
//		isRequest, err := isHttpRequest(ctx, reader)
//		if err != nil {
//			return errors.Wrap(err, `read http request content error`)
//		}
//		if isRequest {
//			req, err := http.ReadRequest(bufReader)
//			if err != nil {
//				return errors.Wrap(err, `read http request content err`)
//			}
//			if err := adapterDump(ctx, copyR, req.Header, req.Body, func() ([]byte, error) {
//				return httputil.DumpRequest(req, false)
//			}); err != nil {
//				return errors.Wrap(err, `dump http request content error`)
//			}
//			return nil
//		}
//
//		isResponse, err := isHttpResponse(ctx, reader)
//		if err != nil {
//			return errors.Wrap(err, `read http response content error`)
//		}
//		if isResponse {
//			resp, err := http.ReadResponse(bufReader, nil)
//			if err != nil {
//				return errors.Wrap(err, `read http response content error`)
//			}
//			if len(resp.TransferEncoding) > 0 && resp.TransferEncoding[0] == "chunked" {
//				chunked, err := http_codec.ReadChunked(bufReader)
//				if err != nil {
//					_ = resp.Body.Close()
//					return errors.Wrap(err, `read http response content error, transfer encoding is chunked`)
//				}
//				_ = resp.Body.Close()
//				buffer := bufutils.NewBufferData(chunked)
//				defer bufutils.ResetBuffer(buffer)
//				resp.Body = ioutil.NopCloser(buffer) // copy
//			}
//			if err := adapterDump(ctx, copyR, resp.Header, resp.Body, func() ([]byte, error) {
//				return httputil.DumpResponse(resp, false)
//			}); err != nil {
//				return errors.Wrap(err, `dump http response content error`)
//			}
//			return nil
//		}
//		return errors.Errorf(`invalid http content`)
//	}
//}

//func adapterDump(ctx context.Context, src *bytes.Buffer, header http.Header, body io.ReadCloser, dumpHeader func() ([]byte, error)) error {
//	defer body.Close()
//	bodyData, err := http_codec.DecodeHttpBody(body, header, false)
//	if err != nil {
//		ctx.Verbose("[HTTP] decode http body err: %v", err)
//		ctx.PrintPayload(src.String())
//		return nil
//	}
//	if len(bodyData) == 0 {
//		ctx.PrintPayload(src.String())
//		return nil
//	}
//	responseHeader, err := dumpHeader()
//	if err != nil {
//		ctx.PrintPayload(src.String())
//		return nil
//	}
//	ctx.PrintPayload(string(responseHeader))
//	ctx.PrintPayload(string(bodyData))
//	return nil
//}
