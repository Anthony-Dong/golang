package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/command/cli"
)

func main() {
	cli.Run(nil, func(ctx context.Context) (*cobra.Command, error) {
		return cli.NewCommand(nil)
	})
}
