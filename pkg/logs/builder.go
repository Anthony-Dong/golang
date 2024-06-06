package logs

import (
	"bytes"
	"context"
	"encoding"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
)

type builder struct {
	level Level
	kvs   []kv
	flag  int
}

type kv struct {
	key   string
	value interface{}
}

func (kv kv) String() string {
	if kv.key == "__str__" {
		return kv.value.(string)
	}
	return fmt.Sprintf("%s=%s", kv.key, toString(kv.value))
}

func (b *builder) Debug() *builder {
	b.level = LevelDebug
	return b
}

func (b *builder) Info() *builder {
	b.level = LevelInfo
	return b
}

func (b *builder) Notice() *builder {
	b.level = LevelNotice
	return b
}

func (b *builder) Warn() *builder {
	b.level = LevelWarn
	return b
}

func (b *builder) Error() *builder {
	b.level = LevelError
	return b
}

func (b *builder) String(format string, v ...interface{}) *builder {
	if len(v) == 0 {
		b.kvs = append(b.kvs, kv{key: "__str__", value: format})
	} else {
		b.kvs = append(b.kvs, kv{key: "__str__", value: fmt.Sprintf(format, v...)})
	}
	return b
}

func (b *builder) KV(key string, value interface{}) *builder {
	b.kvs = append(b.kvs, kv{key: key, value: value})
	return b
}

func (b *builder) Flag(flag int) *builder {
	b.flag = flag
	return b
}

func (b *builder) Emit(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
	// todo add buffer cache
	output := bytes.Buffer{}
	lastIndex := len(b.kvs) - 1
	for index, elem := range b.kvs {
		output.WriteString(elem.String())
		if index == lastIndex {
			break
		}
		output.WriteByte(' ')
	}
	logf(ctx, b.flag, b.level, 2, output.String())
}

func Builder() *builder {
	return &builder{
		flag: defaultLogFlag,
	}
}

func toString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case uint8, uint16, uint32, uint64:
		convertUint64 := func(value interface{}) uint64 {
			switch v := value.(type) {
			case uint8:
				return uint64(v)
			case uint16:
				return uint64(v)
			case uint32:
				return uint64(v)
			case uint64:
				return v
			default:
				panic("ToString uint error")
			}
		}
		return strconv.FormatUint(convertUint64(value), 10)
	case int, int8, int16, int32, int64:
		convertInt64 := func(value interface{}) int64 {
			switch v := value.(type) {
			case int8:
				return int64(v)
			case int16:
				return int64(v)
			case int32:
				return int64(v)
			case int64:
				return v
			case int:
				return int64(v)
			default:
				panic("ToString int error")
			}
		}
		return strconv.FormatInt(convertInt64(value), 10)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case error:
		return v.Error()
	case fmt.Stringer:
		return v.String()
	case []byte:
		return base64.StdEncoding.EncodeToString(v)
	default:
		if v == nil {
			return ""
		}
		if str, isOk := value.(fmt.Stringer); isOk {
			return str.String()
		}
		if codec, isOk := value.(encoding.TextMarshaler); isOk {
			if text, err := codec.MarshalText(); err == nil {
				return string(text)
			}
		}
		if codec, isOk := value.(json.Marshaler); isOk {
			if text, err := codec.MarshalJSON(); err == nil {
				return string(text)
			}
		}
		if result, err := json.Marshal(v); err == nil {
			return string(result)
		}
		return fmt.Sprintf("%v", value)
	}
}
