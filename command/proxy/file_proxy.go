package proxy

import (
	"net/http"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/proxy"
	"github.com/anthony-dong/golang/pkg/utils"
)

func NewFileSystemCommand() (*cobra.Command, error) {
	listenAddr := ""
	fileSystemDir := ""
	cmd := &cobra.Command{
		Use:   "fs",
		Short: "FileSystem Proxy tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			port, err := utils.GetPort(listenAddr)
			if err != nil {
				return err
			}
			ipAddrs, err := utils.GetAllIP(true)
			if err != nil {
				return err
			}
			logs.CtxInfo(ctx, "FileSystem dir: %v", fileSystemDir)
			for _, elem := range ipAddrs {
				logs.CtxInfo(ctx, "FileSystem listen addr: http://%v:%v", elem, port)
			}
			return http.ListenAndServe(listenAddr, proxy.NewFsHandler(fileSystemDir))
		},
	}
	cmd.Flags().StringVarP(&listenAddr, "listen", "l", ":8080", "Listen address")
	cmd.Flags().StringVarP(&fileSystemDir, "dir", "d", utils.GetPwd(), "FileSystem Dir")
	return cmd, nil
}
