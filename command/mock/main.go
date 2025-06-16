package mock

import (
	"fmt"

	"github.com/cloudwego/thriftgo/parser"
	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/idl"
	"github.com/anthony-dong/golang/pkg/rpc"
	"github.com/anthony-dong/golang/pkg/utils"
)

func NewThriftStructMock() (*cobra.Command, error) {
	main := ""
	structName := ""
	cmd := &cobra.Command{
		Use:   "thrift",
		Short: "Generate mock JSON data from Thrift IDL",
		RunE: func(cmd *cobra.Command, args []string) error {
			value, err := GetThriftMockData(main, structName)
			if err != nil {
				return err
			}
			fmt.Println(utils.ToJson(value))
			return nil
		},
	}
	cmd.Flags().StringVar(&structName, "struct", "", "the struct name")
	cmd.Flags().StringVar(&main, "main", "", "the main idl")
	if err := cmd.MarkFlagRequired("main"); err != nil {
		return nil, err
	}
	return cmd, nil
}

func GetThriftMockData(main string, structName string) (interface{}, error) {
	provider := idl.NewMemoryIDLProvider(main)
	thriftIDL, err := provider.ThriftIDL()
	if err != nil {
		return nil, err
	}
	for _, elem := range thriftIDL.Structs {
		if elem.Name == structName {
			thriftType, err := idl.ParseThriftType(thriftIDL, &parser.Type{Name: structName})
			if err != nil {
				return nil, err
			}
			return rpc.GetThriftExampleValue(thriftType, nil, nil)
		}
	}
	return nil, fmt.Errorf(`not found struct [%s] in thrift idl [%s]`, structName, main)
}
