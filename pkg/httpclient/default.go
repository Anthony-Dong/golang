package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func NewCookieAuth(cookie string) Auth {
	return &cookieAuth{
		token: cookie,
	}
}

type cookieAuth struct {
	token string
}

func (c *cookieAuth) Authorization(req *http.Request) error {
	req.Header.Add("Cookie", c.token)
	return nil
}

type defaultBinding struct{}

func (*defaultBinding) Bind(ctx context.Context, respBody []byte, respData interface{}, httpResp *http.Response) (err error) {
	if respData == nil {
		return nil
	}
	switch v := respData.(type) {
	case *string:
		*v = string(respBody)
		return nil
	case *[]byte:
		*v = respBody
		return nil
	}
	if err := json.Unmarshal(respBody, respData); err != nil {
		return fmt.Errorf(`read response body find err. http code: %d, body: '%s'`, httpResp.StatusCode, respBody)
	}
	return nil
}
