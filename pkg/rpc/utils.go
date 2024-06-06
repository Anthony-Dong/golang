package rpc

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/anthony-dong/golang/pkg/utils"
)

func NewRpcRequest(urlStr string, headers []string, body string) (*Request, error) {
	if !strings.Contains(urlStr, "://") {
		urlStr = "thrift://" + urlStr // default is thrift
	}
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	query := parsedUrl.Query()
	req := &Request{
		Protocol:  parsedUrl.Scheme,
		Service:   parsedUrl.Host,
		RPCMethod: strings.TrimPrefix(parsedUrl.Path, "/"),
		Body:      []byte(body),
		Header: utils.MapFromSlice(headers, func(header string) *KV {
			return NewKV(utils.SplitKV(header, ":"))
		}),
		Addr: query.Get("addr"),
		Tag: utils.FlatMapFromMap(query, func(key string, values []string) []*KV {
			if key == "addr" {
				return nil
			}
			ret := make([]*KV, 0, len(values))
			for _, value := range values {
				ret = append(ret, NewKV(key, value))
			}
			return ret
		}),
	}
	if req.Service == "" {
		return nil, fmt.Errorf(`invalid rpc service name. request url like`)
	}
	return req, err
}
