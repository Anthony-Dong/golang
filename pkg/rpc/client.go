package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"
	"strings"
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

func (r *Request) String() string {
	buffer := bytes.Buffer{}
	buffer.WriteString(strings.ToUpper(r.Protocol))
	buffer.WriteString(" ")
	buffer.WriteString(r.Service + "/" + r.RPCMethod)
	if len(r.Tag) > 0 {
		query := url.Values{}
		for _, elem := range r.Tag {
			query.Add(elem.Key, elem.Value)
		}
		buffer.WriteString("?")
		buffer.WriteString(query.Encode())
	}
	buffer.WriteString("\n")
	for _, elem := range r.Header {
		buffer.WriteString(elem.Key)
		buffer.WriteString(": ")
		buffer.WriteString(elem.Value)
		buffer.WriteString("\n")
	}
	if r.Timeout > 0 {
		buffer.WriteString("@timeout")
		buffer.WriteString(": ")
		buffer.WriteString(r.Timeout.String())
		buffer.WriteString("\n")
	}
	if r.EnableModifyRequest {
		buffer.WriteString("@modify: true")
		buffer.WriteString("\n")
	}
	buffer.WriteString("\n")
	if len(r.Body) > 0 {
		buffer.Write(utils.PrettyJsonBytes(r.Body))
	}
	return buffer.String()
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

func (r *Response) String() string {
	buffer := bytes.Buffer{}
	for _, elem := range r.Header {
		buffer.WriteString(elem.Key)
		buffer.WriteString(" ")
		buffer.WriteString(elem.Value)
		buffer.WriteString("\n")
	}
	if r.Spend > 0 {
		buffer.WriteString("@spend")
		buffer.WriteString(": ")
		buffer.WriteString(r.Spend.String())
		buffer.WriteString("\n")
	}
	buffer.WriteString("\n")
	buffer.Write(utils.PrettyJsonBytes(r.Body))
	return buffer.String()
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
