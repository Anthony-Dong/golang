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
		return c.doRequest(ctx, http.MethodGet, path, query, nil, resp)
	}, backoff.WithMaxRetries(backoff.NewConstantBackOff(time.Millisecond*5), 2)); err != nil {
		return err
	}
	return nil
}

func (c *HostClient) Post(ctx context.Context, path string, req interface{}, resp interface{}) error {
	return c.doRequest(ctx, http.MethodPost, path, nil, req, resp)
}

func (c *HostClient) doRequest(ctx context.Context, method string, path string, query map[string]string, req interface{}, resp interface{}) (_err error) {
	if ctx == nil {
		ctx = context.Background()
	}
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
		request.Header.Add("Content-Type", "application/json")
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
	response, err := (&http.Client{Timeout: timeout}).Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

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
