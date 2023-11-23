package rpc

import (
	"context"
	"encoding/json"
	"time"
)

type IDLType string

const (
	IDLTypeLocal IDLType = "local"
)

type Request struct {
	Service  string          `json:"service,omitempty"`
	Method   string          `json:"method,omitempty"`
	Body     json.RawMessage `json:"body,omitempty"`
	Instance Instance        `json:"instance,omitempty"`
	IDLType  IDLType         `json:"idl_type,omitempty"`
	MainIDL  string          `json:"main_idl,omitempty"`
}

type Instance struct {
	Host string `json:"host"`
}

type Response struct {
	Body       json.RawMessage `json:"body,omitempty"`
	TotalSpend time.Duration   `json:"total_spend,omitempty"`
	Spend      time.Duration   `json:"spend,omitempty"`
	Extra      ResponseExtra   `json:"extra,omitempty"`
}

type ResponseExtra struct {
	Instance
	MetaInfo map[string]string `json:"meta_info"`
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
