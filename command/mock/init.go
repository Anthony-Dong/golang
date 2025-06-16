package mock

import (
	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/command"
)

func NewCommand() (*cobra.Command, error) {

	cmd := &cobra.Command{Use: "mock", Short: "Generate mock data for various formats"}

	if err := command.AddCommand(cmd, NewThriftStructMock); err != nil {
		return nil, err
	}

	return cmd, nil
}
