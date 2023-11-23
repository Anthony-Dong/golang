package logs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

var logLevel = LevelInfo
var logFlag = LogFlagPrefix | LogFlagColor | LogFlagTime

const LogFlagPrefix = 1 << 0
const LogFlagColor = 1 << 1
const LogFlagTime = 1 << 2
const LogFlagCaller = 1 << 3

var _levelColor = map[int]func(format string, a ...interface{}) string{
	LevelDebug:  color.HiBlueString,
	LevelInfo:   color.HiCyanString,
	LevelNotice: color.HiGreenString,
	LevelWarn:   color.HiYellowString,
	LevelError:  color.HiRedString,
}

func SetLevel(level int) {
	logLevel = level
}

func SetLevelString(level string) {
	logLevel = stringToLevel(level)
}

func SelFlag(flag int) {
	logFlag = flag
}

func Flush() {

}

const LevelDebug = 0
const LevelInfo = 1
const LevelNotice = 2
const LevelWarn = 3
const LevelError = 4

func stringToLevel(str string) int {
	switch str {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "notice":
		return LevelNotice
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	}
	return logLevel
}

func IsLevel(level int) bool {
	return level >= logLevel
}

func CtxDebug(ctx context.Context, format string, v ...interface{}) {
	logf(ctx, LevelDebug, 2, format, v...)
}

func CtxInfo(ctx context.Context, format string, v ...interface{}) {
	logf(ctx, LevelInfo, 2, format, v...)
}

func CtxWarn(ctx context.Context, format string, v ...interface{}) {
	logf(ctx, LevelWarn, 2, format, v...)
}

func CtxError(ctx context.Context, format string, v ...interface{}) {
	logf(ctx, LevelError, 2, format, v...)
}

func Debug(format string, v ...interface{}) {
	logf(context.Background(), LevelDebug, 2, format, v...)
}

func Info(format string, v ...interface{}) {
	logf(context.Background(), LevelInfo, 2, format, v...)
}

func Notice(format string, v ...interface{}) {
	logf(context.Background(), LevelNotice, 2, format, v...)
}

func Warn(format string, v ...interface{}) {
	logf(context.Background(), LevelWarn, 2, format, v...)
}

func Error(format string, v ...interface{}) {
	logf(context.Background(), LevelError, 2, format, v...)
}

func logf(ctx context.Context, level int, cl int, format string, v ...interface{}) {
	if level < logLevel {
		return
	}
	out := strings.Builder{}

	if logFlag&LogFlagPrefix == LogFlagPrefix {
		switch level {
		case LevelDebug:
			out.WriteString("[DEBUG] ")
		case LevelInfo:
			out.WriteString("[INFO] ")
		case LevelNotice:
			out.WriteString("[NOTICE] ")
		case LevelWarn:
			out.WriteString("[WARN] ")
		case LevelError:
			out.WriteString("[ERROR] ")
		default:
			out.WriteString("[-] ")
		}
	}
	if logFlag&LogFlagTime == LogFlagTime {
		now := time.Now().Format("15:04:05.000")
		out.WriteString(now)
		out.WriteString(" ")
	}

	if logFlag&LogFlagCaller == LogFlagCaller {
		_, file, line, ok := runtime.Caller(cl)
		if !ok {
			file = "???"
			line = 0
		}
		out.WriteString(filepath.Base(file))
		out.WriteString(":")
		out.WriteString(strconv.FormatInt(int64(line), 10))
		out.WriteString(" ")
	}
	out.WriteString(fmt.Sprintf(format, v...))
	output := out.String()
	if logFlag&LogFlagColor == LogFlagColor {
		if foo := _levelColor[level]; foo != nil {
			output = foo(out.String())
		}
	}
	fmt.Println(output)
}

func StdOut(format string, v ...interface{}) {
	fmt.Fprintln(os.Stdout, color.HiCyanString(format, v...))
}

func StdError(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, color.HiRedString(format, v...))
}
