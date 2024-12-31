package strace

import "github.com/spf13/cobra"

func NewCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "strace",
		Short: "strace commands",
		Example: `
1. sudo strace -p <pid> -tt -f -y -s <strsize> -e <expr1,expr2>

- -p <pid>: trace process with process id PID, may be repeated
- -e <expr>: trace, abbrev, verbose, raw, signal, read, write, fault, inject, kvm
- -x: print non-ascii strings in hex
- -xx: print all strings in hex
- -s <strsize>: limit length of print strings to <strsize> chars (default 32)
- -tt: print absolute timestamp with usecs
- -f: follow forks
- -o <file>: send trace output to FILE instead of stderr
- -y: 表示打印出相关文件的路径，对于 open, connect 等函数，它将显示打开或连接的文件或 socket 的路径

2. example:
- sudo strace -p 354481 -tt -f -y -s 65535 -e read,write,close  查看read,write,close三个system call
- sudo strace -p 354481 -tt -f -y -s 65535 查看全部的
`,

		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	return cmd, nil
}
