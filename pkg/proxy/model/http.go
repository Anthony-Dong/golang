package model

import (
	"strings"
	"time"
)

type HTTPCapture struct {
	ID      string    `json:"id"`
	Time    time.Time `json:"time"`
	Spend   int64     `json:"spend"`
	SrcAddr string    `json:"src_addr"`
	IsTLS   bool      `json:"is_tls"`

	Request  *Request  `json:"request"`
	Response *Response `json:"response"`
}

type Request struct {
	Scheme   string        `json:"scheme"`
	Protocol string        `json:"protocol"`
	Method   string        `json:"method"`
	Host     string        `json:"host"`
	Path     string        `json:"path"`
	Header   IgnoreCaseKVS `json:"header"`
	Query    KVS           `json:"query,omitempty"`
	RawBody  []byte        `json:"raw_body,omitempty"`
	Body     string        `json:"body,omitempty"`
}

type Response struct {
	Protocol   string        `json:"protocol"`
	StatusCode int32         `json:"status_code"`
	Header     IgnoreCaseKVS `json:"header"`
	RawBody    []byte        `json:"raw_body,omitempty"`
	Body       string        `json:"body,omitempty"`
}

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type IgnoreCaseKVS []*KV

func (hk IgnoreCaseKVS) Get(key string) string {
	for _, kv := range hk {
		if strings.EqualFold(key, kv.Key) {
			return kv.Value
		}
	}
	return ""
}

type KVS []*KV

func (kvs KVS) Get(key string) string {
	for _, kv := range kvs {
		if key == kv.Key {
			return kv.Value
		}
	}
	return ""
}
