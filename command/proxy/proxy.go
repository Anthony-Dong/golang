package proxy

import (
	"fmt"
	"strings"

	"github.com/anthony-dong/golang/command"

	"github.com/anthony-dong/golang/pkg/proxy/record"

	"github.com/anthony-dong/golang/pkg/logs"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/proxy"
)

func NewCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "proxy",
		Short: "HTTP/HTTPS/Thrift/FileSystem proxy tool",
	}
	if err := command.AddCommand(cmd, NewHTTPCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, NewThriftCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, NewFileSystemCommand); err != nil {
		return nil, err
	}
	return cmd, nil
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

func NewHTTPCommand() (*cobra.Command, error) {
	listenAddr := ""
	cmd := &cobra.Command{
		Use:   "http",
		Short: `HTTP/HTTPS proxy tool`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			port := ""
			if index := strings.LastIndex(listenAddr, ":"); index != -1 {
				port = listenAddr[index+1:]
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
			return proxy.NewProxy(listenAddr, proxy.NewHTTPProxyHandler(proxy.NewRecordHTTPHandler(record.NewConsulStorage()))).Run()
		},
	}
	cmd.Flags().StringVarP(&listenAddr, "listen", "l", ":8080", "Listen address")
	return cmd, nil
}

func NewThriftCommand() (*cobra.Command, error) {
	listenAddr := ""
	dialAddr := ""
	cmd := &cobra.Command{
		Use:   "thrift",
		Short: `Thrift proxy tool`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			return proxy.NewProxy(listenAddr, proxy.NewThriftHandler(dialAddr, record.NewConsulStorage())).Run()
		},
	}
	cmd.Flags().StringVarP(&listenAddr, "listen", "l", ":8080", "Listen address")
	cmd.Flags().StringVar(&dialAddr, "remote", "", "Proxy address")
	return cmd, nil
}
