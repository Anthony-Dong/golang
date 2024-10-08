package rpc

import (
	"context"
	"encoding/json"
	"time"

	"github.com/anthony-dong/golang/pkg/idl"

	"github.com/anthony-dong/golang/pkg/utils"
)

const ProtocolThrift = "thrift"

type Request struct {
	Protocol    string             `json:"protocol,omitempty"`     // thrift/grpc
	ServiceName string             `json:"service_name,omitempty"` // rpc service name
	RPCMethod   string             `json:"rpc_method,omitempty"`   // rpc service method
	Body        json.RawMessage    `json:"body,omitempty"`         // request body
	Header      []*KV              `json:"header,omitempty"`       // request header
	Addr        string             `json:"addr,omitempty"`         // request addr
	Tag         []*KV              `json:"tag,omitempty"`          // custom env/cluster
	Timeout     utils.JsonDuration `json:"timeout,omitempty"`      // timeout
	IDLConfig   *IDLConfig         `json:"idl_config,omitempty"`

	EnableModifyRequest bool `json:"enable_modify_request,omitempty"`
}

func (r *Request) BasicInfo() *Request {
	clone := *r
	clone.Body = nil
	return &clone
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

type IDLConfig struct {
	Main    string   `json:"main,omitempty"`
	Include []string `json:"include,omitempty"`
	Branch  string   `json:"branch,omitempty"`
}

type IDLProvider interface {
	MemoryIDL(ctx context.Context, serviceName string, idlConfig *IDLConfig) (*idl.MemoryIDL, error)
}

func GetValue(kv []*KV, key string) string {
	for _, elem := range kv {
		if elem.Key == key {
			return elem.Value
		}
	}
	return ""
}

func NewKV(key, value string) *KV {
	return &KV{Key: key, Value: value}
}
