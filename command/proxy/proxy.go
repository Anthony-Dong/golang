package proxy

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/proxy"
)

func NewCommand() (*cobra.Command, error) {
	listen := ""
	dial := ""
	command := &cobra.Command{
		Use:   "proxy",
		Short: `Proxy thrift requests`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf(`%v`, r)
				}
			}()
			return proxy.NewProxy(listen, dial, proxy.NewThriftHandler(proxy.ConsoleRecorder)).Run()
		},
	}
	command.Flags().StringVarP(&listen, "listen", "l", "", "The listen addr")
	command.Flags().StringVarP(&dial, "dial", "d", "", "The dial addr")
	if err := command.MarkFlagRequired("listen"); err != nil {
		return nil, err
	}
	if err := command.MarkFlagRequired("dial"); err != nil {
		return nil, err
	}
	return command, nil
}
