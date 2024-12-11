package strace

import (
	"fmt"
	"testing"

	"github.com/anthony-dong/golang/pkg/utils"
)

func TestName(t *testing.T) {
	testString := `arg1, "arg2, with, comma", 'arg3 with \'escaped\' single quotes', arg4,\"arg5 , with \"double\" quotes\",'arg6, with , single quotes'`
	parsedArgs := parseArgs(testString)
	for _, arg := range parsedArgs {
		fmt.Printf("Arg: %s\n", arg)
	}
}

func TestName2(t *testing.T) {
	line, err := parseLogLine(`365094 19:42:35.844286 read(7<socket:[1240417]>, "GET / HTTP/1.1\r\nHost: localhost:8080\r\nUser-Agent: curl/7.64.0\r\nAccept: */*\r\n\r\n", 4096) = 78`)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(utils.ToJson(line, true))

	for _, elem := range line.Args {
		t.Log(hexToString(elem))
	}
}
