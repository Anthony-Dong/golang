package proxy

import (
	"fmt"
	"strings"

	"github.com/anthony-dong/golang/pkg/proxy/record"

	"github.com/anthony-dong/golang/pkg/logs"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/proxy"
)

func NewCommand() (*cobra.Command, error) {
	listen := ""
	dial := ""
	proxyType := ""
	output := ""
	command := &cobra.Command{
		Use:   "proxy",
		Short: `Proxy and Capture thrift/http/https requests`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf(`%v`, r)
				}
			}()
			storage, err := newStorage(output)
			if err != nil {
				return err
			}
			var handler proxy.Handler
			switch proxyType {
			case "http", "https":
				port := ""
				if index := strings.LastIndex(listen, ":"); index != -1 {
					port = listen[index+1:]
				}
				logs.Info(`注意: 安装手册
# 如何配置代理
export http_proxy=http://localhost:%s
export https_proxy=http://localhost:%s
# Linux 如何下载证书 (仅需要操作一次)
1. wget -e http_proxy=localhost:%s http://devtool.mitm/cert/pem -O devtool-ca-cert.pem
2. sudo mv devtool-ca-cert.pem /usr/local/share/ca-certificates/devtool.crt
3. sudo update-ca-certificates
# 注意: 请初始化后再使用才能生效 ....
`, port, port, port)
				handler = proxy.NewHTTPProxyHandler(storage)
			case "thrift":
				handler = proxy.NewThriftHandler(dial, storage)
			default:
				return fmt.Errorf(`invalid proxy type: %s`, proxyType)
			}
			return proxy.NewProxy(listen, handler).Run()
		},
	}
	command.Flags().StringVarP(&listen, "listen", "l", ":8080", "The proxy listen addr.")
	command.Flags().StringVar(&dial, "remote", "", "The remote(thrift) addr.")
	command.Flags().StringVar(&proxyType, "type", "http", "the proxy type. thrift/http/https")
	command.Flags().StringVar(&output, "output", "format", "the output position of the packet (simple/format/json/@file).")
	return command, nil
}

func newStorage(output string) (record.Storage, error) {
	if output == "format" {
		return record.NewConsulStorage(), nil
	}
	if output == "json" {
		return record.NewJsonConsulStorage(), nil
	}
	if output == "simple" {
		return nil, nil
	}
	if strings.HasPrefix(output, "@") {
		return record.NewLocalStorage(strings.TrimPrefix(output, "@"))
	}
	return nil, fmt.Errorf(`invalid output position: %s`, output)
}
