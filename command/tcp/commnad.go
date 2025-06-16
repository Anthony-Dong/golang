package tcp

import (
	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/command"
)

func NewCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{Use: "tcp", Short: "TCP related tools"}
	if err := command.AddCommand(cmd, NewEchoServiceCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, NewHTTPEchoServiceCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, NewBenchmarkEchoServiceCommand); err != nil {
		return nil, err
	}
	return cmd, nil
}
