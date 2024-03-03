package rpc

import (
	"context"
	"encoding/json"
	"time"

	"github.com/anthony-dong/golang/pkg/utils"
)

const ProtocolThrift = "thrift"

type Request struct {
	Protocol  string             `json:"protocol,omitempty"`   // thrift/grpc
	Service   string             `json:"service,omitempty"`    // rpc service name
	RPCMethod string             `json:"rpc_method,omitempty"` // rpc service method
	Body      json.RawMessage    `json:"body,omitempty"`       // request body
	Header    []*KV              `json:"header,omitempty"`     // request header
	Addr      string             `json:"addr,omitempty"`       // request addr
	Tag       []*KV              `json:"tag,omitempty"`        // custom env/cluster
	Timeout   utils.JsonDuration `json:"timeout,omitempty"`    // timeout

	EnableModifyRequest bool `json:"enable_modify_request,omitempty"`
}

func GetValue(kv []*KV, key string) string {
	for _, elem := range kv {
		if elem.Key == key {
			return elem.Value
		}
	}
	return ""
}

type Response struct {
	Body    json.RawMessage `json:"body,omitempty"`
	IsError bool            `json:"is_error,omitempty"`
	Spend   time.Duration   `json:"spend,omitempty"`
	Header  []*KV           `json:"header,omitempty"`
}

type Method struct {
	RPCMethod string `json:"rpc_method,omitempty"`
}

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewKV(key, value string) *KV {
	return &KV{Key: key, Value: value}
}

type IDLInfo struct {
	Main    string   `json:"main"`
	Include []string `json:"include"`
	Branch  string   `json:"branch"`
}

type ExampleCode struct {
	Body json.RawMessage `json:"body"`
}

type Client interface {
	Do(ctx context.Context, req *Request) (*Response, error)
	ListMethods(ctx context.Context) ([]*Method, error)
	GetExampleCode(ctx context.Context, method *Method) (*ExampleCode, error)
}
