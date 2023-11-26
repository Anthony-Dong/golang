package rpc

import (
	"context"
	"encoding/json"
	"time"
)

type Request struct {
	Service  string          `json:"service,omitempty"`  // service name
	Method   string          `json:"method,omitempty"`   // rpc service method
	Body     json.RawMessage `json:"body,omitempty"`     // request body
	Endpoint Endpoint        `json:"endpoint,omitempty"` // request endpoint
}

type Endpoint struct {
	Addr string            `json:"addr,omitempty"`
	Tag  map[string]string `json:"tag,omitempty"` // custom env/cluster
}

type Response struct {
	Body       json.RawMessage `json:"body,omitempty"`
	TotalSpend time.Duration   `json:"total_spend,omitempty"`
	Spend      time.Duration   `json:"spend,omitempty"`
	Extra      ResponseExtra   `json:"extra,omitempty"`
}

type ResponseExtra struct {
	Endpoint
	MetaInfo map[string]string `json:"meta_info,omitempty"`
}

type Method struct {
	RPCMethod string `json:"rpc_method,omitempty"`
	Desc      string `json:"desc,omitempty"`
	Method    string `json:"method,omitempty"`
	Path      string `json:"path,omitempty"`
}

type Client interface {
	Send(ctx context.Context, req *Request) (*Response, error)
	ExampleCode(ctx context.Context, request *Request) (string, error)
	MethodList(ctx context.Context, req *Request) ([]*Method, error)
}
