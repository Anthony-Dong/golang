package proxy

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"io"
	"net"
	"sync"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/anthony-dong/golang/pkg/logs"
)

func NewHTTPProxyHandler(handler fasthttp.RequestHandler) *httpProxyHandler {
	rootCert, err := tls.X509KeyPair(CA_CERT, CA_KEY)
	if err != nil {
		panic(err)
	}
	// cache Leaf
	if rootCert.Leaf, err = x509.ParseCertificate(rootCert.Certificate[0]); err != nil {
		panic(err)
	}
	return &httpProxyHandler{
		RootCert:  rootCert,
		CertCache: &certificateCacheHelper{max: 1024},
		Handler:   handler,
	}
}

type httpProxyHandler struct {
	CertCache *certificateCacheHelper
	RootCert  tls.Certificate
	Handler   fasthttp.RequestHandler

	//EnableMITM func(host string) bool
}

func (t *httpProxyHandler) NewTLSConfig(host string) (*tls.Config, error) {
	config := &tls.Config{InsecureSkipVerify: true}
	certificate, err := NewCertificate(&t.RootCert, []string{stripPort(host)}, t.CertCache)
	if err != nil {
		return nil, err
	}
	config.Certificates = append(config.Certificates, *certificate)
	return config, nil
}

func (t *httpProxyHandler) HandlerConn(conn net.Conn) error {
	httpConfig := NewDefaultHttpConfig()
	handler := t.NewRequestHandler(httpConfig)
	return NewHTTPServer(httpConfig, handler).ServeConn(conn)
}

func (t *httpProxyHandler) NewRequestHandler(cfg *HttpConfig) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		method := string(ctx.Request.Header.Method())
		if method == "CONNECT" {
			defer ctx.Conn().Close()
			if cfg.EnableMITM {
				t.NewMITMProxyRequestHandler(cfg)(ctx)
			} else {
				t.NewTCPProxyRequestHandler()(ctx)
			}
			return
		}
		t.NewHostProxyRequestHandler(false)(ctx)
	}
}

func safeLoadHost(ctx *fasthttp.RequestCtx) string {
	host := string(ctx.Request.Host())
	if host == "" {
		// CONNECT search.maven.org:443 HTTP/1.0
		host = string(ctx.Request.Header.RequestURI())
	}
	return host
}

// HTTPS代理服务器的工作原理基于HTTP的CONNECT方法。与普通的HTTP代理（根据客户端的GET，POST等请求，代理服务器会直接进行处理并将结果返回给客户端）有所不同，HTTPS代理主要用于建立一个TCP的隧道，用于客户端和目标服务器的相互通信。
// 1. 客户端向代理服务器发送一个CONNECT请求，请求中包含目标服务器的地址和端口号。
// 2. 如果代理服务器允许此连接，它会与目标服务器建立TCP连接，并向客户端发送一个状态行，比如 "HTTP/1.0 200 Connection Established"，告知客户端可以开始发送请求到目标服务器了。
// 3. 此时，代理服务器主要扮演数据转发的角色，客户端与目标服务器之间可以通过这个隧道来发送任何类型的数据，包括但不限于HTTP请求。
func (t *httpProxyHandler) NewTCPProxyRequestHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		host := string(ctx.Request.Host())
		srcConn := ctx.Conn()
		dstConn, err := net.Dial("tcp", host)
		if err != nil {
			returnError(ctx, err, `Dst Connection Establish Error`)
			return
		}
		defer dstConn.Close()
		if _, err := srcConn.Write([]byte("HTTP/1.0 200 Connection Established\r\n\r\n")); err != nil {
			returnError(ctx, err, `Connection Establish Error`)
			return
		}

		logs.CtxDebug(ctx, "Establish %s -> %s tcp conn", srcConn.RemoteAddr(), dstConn.RemoteAddr())
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			if _, err := io.Copy(srcConn, dstConn); err != nil {
				logs.CtxError(ctx, "transfer %s -> %s tcp conn find err: %v", srcConn.RemoteAddr(), dstConn.RemoteAddr(), err)
			}
		}()
		go func() {
			defer wg.Done()
			if _, err := io.Copy(dstConn, srcConn); err != nil {
				logs.CtxError(ctx, "transfer %s -> %s tcp conn find err: %v", dstConn.RemoteAddr(), srcConn.RemoteAddr(), err)
			}
		}()
		wg.Wait()
		setNoResponse(ctx)
	}
}

// NewMITMProxyRequestHandler MITM(Man-In-The-Middle)中间人攻击
func (t *httpProxyHandler) NewMITMProxyRequestHandler(cfg *HttpConfig) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		host := string(ctx.Request.Host())
		if host == "" {
			// CONNECT search.maven.org:443 HTTP/1.0
			host = string(ctx.Request.Header.RequestURI())
		}
		srcConn := ctx.Conn()
		if _, err := srcConn.Write([]byte("HTTP/1.0 200 Connection Established\r\n\r\n")); err != nil {
			returnError(ctx, err, `Connection Establish Error`)
			return
		}
		logs.CtxDebug(ctx, "Establish %s -> %s conn success", srcConn.RemoteAddr(), host)
		tlsConfig, err := t.NewTLSConfig(host)
		if err != nil {
			returnError(ctx, err, `Create TLS Config Error`)
			return
		}
		if err := NewHTTPServer(cfg, t.NewHostProxyRequestHandler(true)).ServeConn(tls.Server(srcConn, tlsConfig)); err != nil {
			returnError(ctx, err, `Handler TLS Connection Error`)
			return
		}
		setNoResponse(ctx)
	}
}

func returnError(ctx *fasthttp.RequestCtx, err error, msg string) {
	ctx.Error(msg+": "+err.Error(), fasthttp.StatusInternalServerError)
	logs.CtxError(ctx, "return err: %v, msg: %s", err, msg)
}

func setNoResponse(ctx *fasthttp.RequestCtx) {
	ctx.HijackSetNoResponse(true)
	ctx.Hijack(func(_ net.Conn) {})
}

func (t *httpProxyHandler) NewHostProxyRequestHandler(isTls bool) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if t.isDownloadCertPem(ctx) {
			return
		}
		ctx.SetUserValue("is_tls", isTls)
		t.Handler(ctx)
	}
}

func (t *httpProxyHandler) isDownloadCertPem(ctx *fasthttp.RequestCtx) bool {
	req := &ctx.Request
	if bytes.Contains(req.Host(), []byte("devtool.mitm")) && bytes.Equal(req.RequestURI(), []byte("/cert/pem")) {
		ctx.Response.SetBody(CA_CERT)
		ctx.Response.Header.Set("Content-Type", "application/x-x509-ca-cert")
		ctx.Response.Header.Set("Content-Disposition", "attachment; filename=devtool-ca-cert.pem")
		return true
	}
	return false
}

func NewHTTPServer(cfg *HttpConfig, handler fasthttp.RequestHandler) *fasthttp.Server {
	return &fasthttp.Server{
		StreamRequestBody:  cfg.StreamRequestBody,
		ReadBufferSize:     cfg.ReadBufferSize,
		ReadTimeout:        cfg.ReadTimeout,
		WriteTimeout:       cfg.WriteTimeout,
		IdleTimeout:        cfg.IdleTimeout,
		MaxRequestBodySize: cfg.MaxRequestBodySize,

		// multipart/form-data 使用 stream 处理，防止大文件把内存给撑爆了
		DisablePreParseMultipartForm: true,

		// skip server、data、content-type header.
		NoDefaultServerHeader: true,
		NoDefaultDate:         true,
		NoDefaultContentType:  true,

		// 处理客户端的预检请求
		ContinueHandler: func(header *fasthttp.RequestHeader) bool { return true },
		Handler:         handler,
		CloseOnShutdown: true,
	}
}

func NewHTTPClient(cfg *HttpConfig) *fasthttp.Client {
	return &fasthttp.Client{
		ReadTimeout:              cfg.ReadTimeout,
		WriteTimeout:             cfg.WriteTimeout,
		MaxResponseBodySize:      cfg.MaxResponseBodySize,
		ReadBufferSize:           cfg.ReadBufferSize,
		StreamResponseBody:       cfg.StreamResponseBody,
		NoDefaultUserAgentHeader: true,
		DialDualStack:            true, // dns解析支持ipv6/ipv4
		DisablePathNormalizing:   true, // 禁止对path进行规范化处理
	}
}

type HttpConfig struct {
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	IdleTimeout         time.Duration // is http keepalive time
	MaxRequestBodySize  int
	MaxResponseBodySize int
	ReadBufferSize      int
	StreamRequestBody   bool
	StreamResponseBody  bool

	EnableMITM bool
}

func NewDefaultHttpConfig() *HttpConfig {
	return &HttpConfig{
		ReadTimeout:         15 * time.Minute,
		WriteTimeout:        time.Duration(0),
		IdleTimeout:         15 * time.Minute,
		MaxRequestBodySize:  2 * (1 << 20), // 2m
		MaxResponseBodySize: 2 * (1 << 20), // 2m
		ReadBufferSize:      1 * (1 << 20), // 1m

		// stream body 的原理是仅读取(MaxRequestBodySize/MaxResponseBodySize)的body到内存中，超过的通过stream读取
		StreamRequestBody:  true,
		StreamResponseBody: true,

		// 开启MITM，启动https中间人攻击
		EnableMITM: true,
	}
}
