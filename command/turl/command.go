package turl

import (
	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/idl"
	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/rpc"
	"github.com/anthony-dong/golang/pkg/utils"
)

// NewTurlCommand
// turl --idl ./api.thrift --method RPCAPI1 --data '{}'  --addr localhost:8888
func NewTurlCommand() (*cobra.Command, error) {
	req := &rpc.Request{}
	tags := make([]string, 0)
	listMethod := false
	genExampleCode := false
	mainIDL := ""
	reqBody := ""
	cmd := &cobra.Command{
		Use:   "turl",
		Short: "Send thrift request like curl",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			client := rpc.NewThriftClient(idl.NewLocalIDLProvider(mainIDL))
			for _, tag := range tags {
				k, v := utils.ReadKVByColon(tag)
				if k != "" {
					req.Endpoint.Tag[k] = v
				}
			}
			if listMethod {
				list, err := client.MethodList(ctx, req)
				if err != nil {
					return err
				}
				logs.StdOut(utils.ToJson(list, true))
				return nil
			}
			if genExampleCode {
				example, err := client.ExampleCode(ctx, req)
				if err != nil {
					return err
				}
				logs.StdOut(utils.PrettyJson(example))
				return nil
			}
			logs.CtxInfo(ctx, "req url: thrift://%s/%s", req.Endpoint.Addr, req.Method)
			req.Body = []byte(reqBody)
			logs.CtxInfo(ctx, "req body:\n%s", req.Body)
			resp, err := client.Send(ctx, req)
			if err != nil {
				return err
			}
			logs.CtxInfo(ctx, "req spend: %s", utils.ToString(resp.Spend))
			logs.CtxInfo(ctx, "resp body:\n%s", utils.PrettyJson(string(resp.Body)))
			return nil
		},
	}
	cmd.Flags().StringVar(&mainIDL, "idl", "", "The service main idl")
	cmd.Flags().StringVar(&req.Service, "service", "", "The service name")
	cmd.Flags().StringVar(&req.Method, "method", "", "The request method name")
	cmd.Flags().StringVar(&reqBody, "data", "", "The request data")
	cmd.Flags().StringVar(&req.Endpoint.Addr, "addr", "", "The request endpoint addr")
	cmd.Flags().StringSliceVar(&tags, "tag", []string{}, "The request endpoint tag")
	cmd.Flags().BoolVar(&listMethod, "methods", false, "List methods")
	cmd.Flags().BoolVar(&genExampleCode, "example", false, "Gen example code")
	if err := cmd.MarkFlagRequired("idl"); err != nil {
		return nil, err
	}
	return cmd, nil
}
