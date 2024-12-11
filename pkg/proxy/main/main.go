package main

import (
	"github.com/anthony-dong/golang/pkg/proxy"
	"github.com/anthony-dong/golang/pkg/proxy/record"
)

func main() {
	if err := proxy.NewProxy(":8080", proxy.NewHTTPProxyHandler(proxy.NewRecordHTTPHandler(record.NewConsulStorage()))).Run(); err != nil {
		panic(err)
	}
}
