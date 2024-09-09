package dag

import (
	"bytes"
	"fmt"
	"sync"
)

type Trace interface {
	Report(from string, to string, option map[string]string)
}

type defaultTrace struct {
	lock    sync.Mutex
	reports []struct {
		From   string
		To     string
		Option map[string]string
	}
}

func NewTrace() *defaultTrace {
	return &defaultTrace{}
}

func (t *defaultTrace) Report(from string, to string, option map[string]string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.reports = append(t.reports, struct {
		From   string
		To     string
		Option map[string]string
	}{From: from, To: to, Option: option})
}

func (t *defaultTrace) Graphviz() string {
	builder := bytes.NewBuffer(nil)
	builder.WriteString("digraph G {\n")
	for _, report := range t.reports {
		if report.From == "" {
			builder.WriteString(fmt.Sprintf("  %q ;\n", report.To))
			continue
		}
		builder.WriteString(fmt.Sprintf("  %q -> %q ;\n", report.From, report.To))
	}
	builder.WriteString("}\n")
	return builder.String()
}
