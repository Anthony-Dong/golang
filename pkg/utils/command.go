package utils

import (
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"golang.org/x/term"

	"github.com/anthony-dong/golang/pkg/logs"
)

// RunCmdWithShell 会启动一个shell终端帮助命令执行! 方便进入容器内进行调试！
func RunCmdWithShell(cmd *exec.Cmd) error {
	// Start the command with a pty.
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return err
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH) // SIGWINCH is 窗口大小改变的信号.
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH                        // Initial resize.
	defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.

	// Set stdin in raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	// NOTE: The goroutine will keep reading until the next keystroke before returning.
	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, ptmx)

	return nil
}

func RunCmd(cmd *exec.Cmd, logPrefix string) error {
	if _, err := RunDaemonCmd(cmd, logPrefix, false); err != nil {
		return err
	}
	return nil
}

func RunCommand(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RunDaemonCmd(cmd *exec.Cmd, logPrefix string, isDaemon bool) (func(), error) {
	logs.Info("%scmd name: %s", logPrefix, cmd.Path)
	if cmd.Dir == "" {
		dir, err := os.Getwd()
		if err == nil {
			logs.Info("%scmd work dir: %s", logPrefix, dir)
		}
	} else {
		logs.Info("%scmd work dir: %s", logPrefix, cmd.Dir)
	}
	logs.StdOut(logPrefix + "cmd exec: \n" + PrettyCmd(cmd, 0))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	logs.Info("%scmd start success. pid: %d", logPrefix, cmd.Process.Pid)
	if !isDaemon {
		logs.Info("%scmd waiting. pid: %d", logPrefix, cmd.Process.Pid)
		return nil, cmd.Wait()
	}
	return func() {
		logs.Info("%scmd waiting. pid: %d", logPrefix, cmd.Process.Pid)
		if err := cmd.Wait(); err != nil {
			logs.Info("%scmd waiting find err: %v. pid: %d", err)
		}
	}, nil
}
