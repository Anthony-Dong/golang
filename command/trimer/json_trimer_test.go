package trimer

import (
	"testing"
)

func TestTrimJsonNUll(t *testing.T) {
	out, err := TrimJson([]byte(`{
    "PackScene": 161,
    "BizCode": "food",
    "magic": 2,
    "Base": {
        "LogID": "",
        "Caller": "",
        "Addr": "",
        "Client": "",
        "TrafficEnv": {
            "Open": false,
            "Env": ""
        },
        "Extra": {
            "": "",
            "env": "ppe_discount_upgrade"
        }
    }
}`))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(out))
}
