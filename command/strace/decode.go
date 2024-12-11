package strace

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ReturnType string

const TypeResumed ReturnType = "resumed"
const TypeUnfinished ReturnType = "unfinished"

type Trace struct {
	PID        int
	Time       time.Time
	Func       string
	Args       []string
	Return     int
	ReturnType ReturnType
}

func (*Trace) GetSocket() {

}

var traceLogRegexp = regexp.MustCompile(`(\d+)\s+([\d:.]+)\s+(\w+)`)

func parseLogLine(logLine string) (*Trace, error) {
	matches := traceLogRegexp.FindStringSubmatch(logLine)
	if len(matches) < 6 {
		return nil, fmt.Errorf("log line does not match expected format")
	}
	pid, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("invalid PID: %v", err)
	}
	timeStr := matches[2]
	currentTime := time.Now()
	logTime, err := time.ParseInLocation("15:04:05.000000", timeStr, currentTime.Location())
	if err != nil {
		return nil, fmt.Errorf("invalid time format: %v", err)
	}
	// Set the date to today; only the time part is parsed from the log
	logTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), logTime.Hour(), logTime.Minute(), logTime.Second(), logTime.Nanosecond(), logTime.Location())

	function := matches[3]
	argsStr := matches[4]
	returnStr := matches[5]
	returnValue := 0
	if returnStr != "" {
		atoi, err := strconv.Atoi(returnStr)
		if err != nil {
			return nil, fmt.Errorf("invalid return: %v", err)
		}
		returnValue = atoi
	}
	args := parseArgs(argsStr)

	trace := &Trace{
		PID:    pid,
		Time:   logTime,
		Func:   function,
		Args:   args,
		Return: returnValue,
	}
	return trace, nil
}

func parseArgs(argsStr string) []string {
	var args []string
	arg := ""
	inDoubleQuotes := false
	inSingleQuotes := false
	isEscaped := false

	hasNext := func(index int) bool {
		return index+1 <= len(argsStr)
	}

	for index, r := range argsStr {
		switch r {
		case '\\':
			if isEscaped {
				arg += string(r)
				isEscaped = false
			} else if hasNext(index) {
				if inDoubleQuotes && argsStr[index+1] == '"' {
					isEscaped = true
				} else if inSingleQuotes && argsStr[index+1] == '\'' {
					isEscaped = true
				} else {
					arg += string(r)
				}
			} else {
				arg += string(r)
			}
		case '"':
			if isEscaped {
				arg += string(r)
				isEscaped = false
			} else if inSingleQuotes {
				arg += string(r)
			} else {
				inDoubleQuotes = !inDoubleQuotes
				//arg += string(r)
			}
		case '\'':
			if isEscaped {
				arg += string(r)
				isEscaped = false
			} else if inDoubleQuotes {
				arg += string(r)
			} else {
				inSingleQuotes = !inSingleQuotes
				//arg += string(r)
			}
		case ',':
			if isEscaped || inDoubleQuotes || inSingleQuotes {
				arg += string(r)
				isEscaped = false
			} else {
				args = append(args, strings.TrimSpace(arg))
				arg = ""
			}
		default:
			arg += string(r)
			isEscaped = false
		}
	}
	if arg != "" {
		args = append(args, strings.TrimSpace(arg))
	}
	return args
}

func hexToString(input string) string {
	decodeString, err := hex.DecodeString(strings.ReplaceAll(input, "\\x", ""))
	if err != nil {
		return input
	}
	return string(decodeString)
}
