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
		Protocol:    parsedUrl.Scheme,
		ServiceName: parsedUrl.Host,
		RPCMethod:   strings.TrimPrefix(parsedUrl.Path, "/"),
		Body:        []byte(body),
		Header: utils.MapFromSlice(headers, func(header string) *KV {
			return NewKV(utils.SplitKV(header, ":"))
		}),
		Addr: query.Get("addr"),
		Tag:  GetTagFromQuery(query, []string{"addr"}),
	}
	if req.ServiceName == "" {
		return nil, fmt.Errorf(`invalid rpc service name. request url like`)
	}
	return req, err
}

func GetTagFromQuery(query url.Values, slip []string) []*KV {
	return utils.FlatMapFromMap(query, func(key string, values []string) []*KV {
		if utils.Contains(slip, key) {
			return nil
		}
		ret := make([]*KV, 0, len(values))
		for _, value := range values {
			ret = append(ret, NewKV(key, value))
		}
		return ret
	})
}
