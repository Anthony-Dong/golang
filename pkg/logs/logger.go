package logs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/fatih/color"
)

const LogIDKey = "LOG_ID"

var defaultLogLevel Level = LevelInfo
var defaultLogFlag = LogFlagPrefix | LogFlagColor | LogFlagTime | LogFlagLogID

var writer io.Writer = os.Stdout

func SetWriter(w io.Writer) {
	writer = w
}

func SetFileWriter(filename string) {
	open, err := os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	writer = open
}

func SetPrinterStdError() {
	writer = os.Stderr
}

const LogFlagPrefix = 1 << 0
const LogFlagColor = 1 << 1
const LogFlagTime = 1 << 2
const LogFlagCaller = 1 << 3
const LogFlagLogID = 1 << 4

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

func GetLogId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	value, _ := ctx.Value(LogIDKey).(string)
	return value
}

func CtxWithLogID(ctx context.Context, id string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, LogIDKey, id)
}

func logf(ctx context.Context, flag int, level Level, cl int, format string, v ...interface{}) {
	if level < defaultLogLevel {
		return
	}
	if writer == nil {
		return
	}
	out := bytes.Buffer{}
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

	if flag&LogFlagLogID == LogFlagLogID {
		if logId := GetLogId(ctx); logId != "" {
			out.WriteByte('[')
			out.WriteString(logId)
			out.WriteString("] ")
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
	payload := format
	if len(v) > 0 {
		payload = fmt.Sprintf(format, v...)
	}
	out.WriteString(payload)
	out.WriteByte('\n')
	if !color.NoColor && flag&LogFlagColor == LogFlagColor {
		if colorFunc := _levelColor[level]; colorFunc != nil {
			colorOut := colorFunc(out.String())
			out.Reset()
			out.WriteString(colorOut)
		}
	}
	_, err := writer.Write(out.Bytes())
	if err != nil {
		log.Printf("writing log error: %v", err)
	}
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
