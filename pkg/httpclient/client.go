package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/anthony-dong/golang/pkg/logs"

	"github.com/cenkalti/backoff"
)

const defaultTimeout = time.Second * 60

type HostClient struct {
	Timeout time.Duration
	Auth    Auth
	Binding Binding
	Host    string
}

type Auth interface {
	Authorization(req *http.Request) error
}

type Binding interface {
	Bind(ctx context.Context, respBody []byte, respData interface{}, httpResp *http.Response) (err error)
}

func (c *HostClient) Get(ctx context.Context, path string, query map[string]string, resp interface{}) error {
	if err := backoff.Retry(func() error {
		return c.doRequest(ctx, http.MethodGet, path, query, nil, nil, resp)
	}, backoff.WithMaxRetries(backoff.NewConstantBackOff(time.Millisecond*5), 2)); err != nil {
		return err
	}
	return nil
}

func (c *HostClient) Post(ctx context.Context, path string, req interface{}, resp interface{}) error {
	return c.doRequest(ctx, http.MethodPost, path, nil, nil, req, resp)
}

func (c *HostClient) Put(ctx context.Context, path string, req interface{}, resp interface{}) error {
	return c.doRequest(ctx, http.MethodPost, path, nil, nil, req, resp)
}

func (c *HostClient) Delete(ctx context.Context, path string, query map[string]string, resp interface{}) error {
	return c.doRequest(ctx, http.MethodDelete, path, query, nil, nil, resp)
}

func (c *HostClient) Do(ctx context.Context, method string, path string, query map[string]string, header map[string]string, req interface{}, resp interface{}) error {
	return c.doRequest(ctx, method, path, query, header, req, resp)
}

func (c *HostClient) doRequest(ctx context.Context, method string, path string, query map[string]string, header map[string]string, req interface{}, resp interface{}) (_err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	start := time.Now()
	log := logs.Builder().Info().Prefix("http request")
	defer func() {
		log.KV("spend", fmt.Sprintf("%dms", time.Now().Sub(start)/time.Millisecond))
		if _err != nil {
			log.Error()
			log.KV("err", fmt.Sprintf(`{%s}`, _err.Error()))
		}
		log.Emit(ctx)
	}()
	var reqJson bool
	var reqBody io.Reader
	if req != nil {
		reqBodyData, err := json.Marshal(req)
		if err != nil {
			return err
		}
		reqBody = bytes.NewBuffer(reqBodyData)
		reqJson = true
	}
	request, err := http.NewRequest(method, c.Host+path, reqBody)
	if err != nil {
		return backoff.Permanent(fmt.Errorf(`new http request find err: %v`, err))
	}
	if request.Header == nil {
		request.Header = map[string][]string{}
	}
	if reqJson {
		request.Header.Set("Content-Type", "application/json")
	}
	for k, v := range header {
		request.Header.Set(k, v)
	}
	if c.Auth != nil {
		if err := c.Auth.Authorization(request); err != nil {
			return err
		}
	}
	queryV := request.URL.Query()
	for k, v := range query {
		queryV.Set(k, v)
	}
	request.URL.RawQuery = queryV.Encode()
	timeout := c.Timeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	log.KV("method", request.Method)
	log.KV("host", fmt.Sprintf("%s://%s", request.URL.Scheme, request.URL.Host))
	log.KV("path", request.URL.Path)
	log.KV("query", fmt.Sprintf(`{%s}`, request.URL.RawQuery))
	response, err := (&http.Client{Timeout: timeout}).Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	log.KV("status_code", response.StatusCode)
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf(`read response body find err: %v. http code: %d`, err, response.StatusCode)
	}
	if c.Binding == nil {
		c.Binding = &defaultBinding{}
	}
	if err := c.Binding.Bind(ctx, respBody, resp, response); err != nil {
		return err
	}
	return nil
}
