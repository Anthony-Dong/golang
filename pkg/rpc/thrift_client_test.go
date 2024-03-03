package rpc

import (
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

func TestNewRequest(t *testing.T) {
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
