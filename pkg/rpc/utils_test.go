package rpc

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/anthony-dong/golang/pkg/utils"
)

func TestToRpcRequest(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "thrift://xxx.xxx.xxx/RPCMethod?env=xxx", bytes.NewBufferString(`{"k1": "v1"}`))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("log_id", "xxxxxxxxx")
	rpcRequest, err := ToRpcRequest(request)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(utils.ToJson(rpcRequest, true))
}
