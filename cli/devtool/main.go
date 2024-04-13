package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/anthony-dong/golang/pkg/logs"

	"github.com/anthony-dong/golang/command"
	"github.com/anthony-dong/golang/command/cli"
)

func main() {
	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)
	defer command.Close()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		cmd, err := cli.NewCommand(nil)
		if err != nil {
			command.Defer(func() { command.ExitError(err) })
			return
		}
		if err := cmd.ExecuteContext(ctx); err != nil {
			command.Defer(func() { command.ExitError(err) })
			return
		}
	}()
	select {
	case v, _ := <-done:
		logs.Debug("process (%d) receive signal (%s) done", os.Getpid(), v)
	case <-ctx.Done():
		logs.Debug("process (%d) done", os.Getpid())
	}
}
