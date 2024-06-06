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
	"unicode"

	"github.com/fatih/color"
)

var defaultLogLevel Level = LevelInfo
var defaultLogFlag = LogFlagPrefix | LogFlagColor | LogFlagTime
var print func(output string) = func(output string) {
	fmt.Print(output)
}

const LogFlagPrefix = 1 << 0
const LogFlagColor = 1 << 1
const LogFlagTime = 1 << 2
const LogFlagCaller = 1 << 3

var _levelColor = map[Level]func(format string, a ...interface{}) string{
	LevelDebug:  color.HiBlueString,
	LevelInfo:   color.HiCyanString,
	LevelNotice: color.HiGreenString,
	LevelWarn:   color.HiYellowString,
	LevelError:  color.HiRedString,
}

func SetLevel(level Level) {
	defaultLogLevel = level
}

func SetPrinter(printer func(output string)) {
	print = printer
}

func SetPrinterStdError() {
	print = func(output string) {
		fmt.Fprint(os.Stderr, output)
	}
}

func SetLevelString(level string) {
	ll, isExist := stringLevelMap[level]
	if !isExist {
		return
	}
	defaultLogLevel = ll
}

func SelFlag(flag int) {
	defaultLogFlag = flag
}

func LogLevel() Level {
	return defaultLogLevel
}

func Flush() {

}

type Level uint8

const (
	LevelDebug Level = iota
	LevelInfo
	LevelNotice
	LevelWarn
	LevelError
)

var (
	levelStringMap = map[Level]string{
		LevelDebug:  "debug",
		LevelInfo:   "info",
		LevelNotice: "notice",
		LevelWarn:   "warn",
		LevelError:  "error",
	}
	stringLevelMap = map[string]Level{}
)

func init() {
	for level, str := range levelStringMap {
		stringLevelMap[str] = level
	}
}

func (l Level) String() string {
	str, isExist := levelStringMap[l]
	if isExist {
		return str
	}
	return "level-" + strconv.Itoa(int(l))
}

func IsLevel(level Level) bool {
	return level >= defaultLogLevel
}

func CtxDebug(ctx context.Context, format string, v ...interface{}) {
	logf(ctx, defaultLogFlag, LevelDebug, 2, format, v...)
}

func CtxInfo(ctx context.Context, format string, v ...interface{}) {
	logf(ctx, defaultLogFlag, LevelInfo, 2, format, v...)
}

func CtxWarn(ctx context.Context, format string, v ...interface{}) {
	logf(ctx, defaultLogFlag, LevelWarn, 2, format, v...)
}

func CtxError(ctx context.Context, format string, v ...interface{}) {
	logf(ctx, defaultLogFlag, LevelError, 2, format, v...)
}

func Debug(format string, v ...interface{}) {
	logf(context.Background(), defaultLogFlag, LevelDebug, 2, format, v...)
}

func IsDebug() bool {
	return IsLevel(LevelDebug)
}

func Info(format string, v ...interface{}) {
	logf(context.Background(), defaultLogFlag, LevelInfo, 2, format, v...)
}

func Notice(format string, v ...interface{}) {
	logf(context.Background(), defaultLogFlag, LevelNotice, 2, format, v...)
}

func Warn(format string, v ...interface{}) {
	logf(context.Background(), defaultLogFlag, LevelWarn, 2, format, v...)
}

func Error(format string, v ...interface{}) {
	logf(context.Background(), defaultLogFlag, LevelError, 2, format, v...)
}

func logf(_ context.Context, flag int, level Level, cl int, format string, v ...interface{}) {
	if level < defaultLogLevel {
		return
	}
	if print == nil {
		return
	}
	out := strings.Builder{}

	if flag&LogFlagPrefix == LogFlagPrefix {
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
	if flag&LogFlagTime == LogFlagTime {
		now := time.Now().Format("15:04:05.000")
		out.WriteString(now)
		out.WriteString(" ")
	}

	if flag&LogFlagCaller == LogFlagCaller {
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
	logData := fmt.Sprintf(format, v...)
	logData = trimRightSpace(logData)
	out.WriteString(logData)
	out.WriteByte('\n')
	output := out.String()
	if flag&LogFlagColor == LogFlagColor {
		if foo := _levelColor[level]; foo != nil {
			output = foo(out.String())
		}
	}
	print(output)
}

func StdOut(format string, v ...interface{}) {
	fmt.Fprintln(os.Stdout, color.HiCyanString(format, v...))
}

func StdError(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, color.HiRedString(format, v...))
}

func trimRightSpace(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.TrimRightFunc(str, func(r rune) bool {
		return unicode.IsSpace(r)
	})
}
