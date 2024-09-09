package utils

import (
	"os"
	"os/exec"

	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/utils/internal"
)

// RunCmdWithPty 会启动一个shell终端帮助命令执行! 方便进入容器内进行调试！
func RunCmdWithPty(cmd *exec.Cmd) error {
	return internal.RunPty(cmd)
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
			logs.Info("%scmd waiting find err: %v. pid: %d", logPrefix, err, cmd.Process.Pid)
		}
	}, nil
}
