package logs

import (
	"bytes"
	"context"
	"fmt"
)

type builder struct {
	level  Level
	prefix string
	kvs    []kv
}

type kv struct {
	key   string
	value interface{}
}

func (kv kv) String() string {
	if kv.key == "__str__" {
		return kv.value.(string)
	}
	return fmt.Sprintf("%s=%v", kv.key, kv.value)
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

func (b *builder) Emit(ctx context.Context) {
	output := bytes.Buffer{}
	for index, elem := range b.kvs {
		output.WriteString(elem.String())
		if index == len(b.kvs)-1 {
			continue
		}
		output.WriteString(" ")
	}
	logf(ctx, b.level, 2, output.String())
}

func Builder() *builder {
	return &builder{}
}
