package rpc

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"testing"
)

func Test_ThriftClient(t *testing.T) {
	return
	// local test
	/*mainIdl := filepath.Join(utils.GetGoProjectDir(), "pkg/idl/test/api.thrift")
	provider := idl.NewDescriptorProvider(idl.NewMemoryIDLProvider(mainIdl))
	client := NewThriftClient(provider)
	ctx := context.Background()
	req := &Request{
		Service:   "TestServiceName",
		RPCMethod: "RPCAPI1",
		Body:      []byte(`{"Field1": "success"}`),
		Addr:      "127.0.0.1:8888",
	}
	jsonExample, err := NewThriftJsonExample(ctx, provider, "RPCAPI1")
	t.Logf("example code: %s\n", jsonExample)

	req.Body = []byte(jsonExample)
	send, err := client.Do(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(utils.ToJson(send, true))*/
}

func TestName(t *testing.T) {
	parse, err := url.Parse(`thrift://xxx.xxx.xxx/RPCMethod?env=xxx`)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(parse)
	request, err := http.NewRequest("", "thrift://xxx.xxx.xxx/RPCMethod?env=xxx", nil)
	if err != nil {
		t.Fatal(err)
	}
	dumpRequest, err := httputil.DumpRequest(request, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(dumpRequest))
}

func TestThriftClient_InjectBaseRequestInfo(t *testing.T) {
	client := &ThriftClient{}
	result, err := client.InjectBaseRequestInfo(context.Background(), []byte(`{"Field1": "success", "Base": {"Extra": {"user_extra": {"k1": "v1"}}}}`), map[string]string{
		"k2": "v2",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(result))
}
