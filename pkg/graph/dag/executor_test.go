package dag

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExecutor_Simple(t *testing.T) {
	graph := NewGraph()
	start := time.Now()
	graph.MustNode("a", func(ctx context.Context) {
		time.Sleep(time.Millisecond * 10)
		info(t, "a")
	})

	graph.MustNode("b", func(ctx context.Context) {
		time.Sleep(time.Millisecond * 10)
		info(t, "b")
	})

	graph.MustNode("c", func(ctx context.Context) error {
		time.Sleep(time.Millisecond * 10)
		info(t, "c")
		return nil
	})
	if err := NewNodeDependencyBuilder(graph).Node("a").To("b", "c").Build(); err != nil {
		t.Fatal(err)
	}
	execGraph(t, graph)
	info(t, "success")
	spend := time.Since(start)
	assert.Equal(t, spend >= time.Millisecond*20, true)
	assert.Equal(t, spend < time.Millisecond*25, true)

}

func execGraph(t testing.TB, graph *Graph) {
	executor := NewExecutor()
	if err := executor.Execute(context.Background(), graph); err != nil {
		t.Fatal(err)
	}
}

func info(t testing.TB, format string, args ...interface{}) {
	t.Logf(time.Now().Format("15:04:05.000")+": "+format+"\n", args...)
}

func TestExecutor_Cond(t *testing.T) {
	graph := NewGraph()
	name := ""
	graph.MustNode("a", func(ctx context.Context) (string, error) {
		return "b", nil
	})
	graph.MustNode("b", func(ctx context.Context) {
		name = ContextCurNode(ctx)
	})
	graph.MustNode("c", func(ctx context.Context) {
		name = ContextCurNode(ctx)
	})
	if err := NewNodeDependencyBuilder(graph).Node("a").To("b", "c").Build(); err != nil {
		t.Fatal(err)
	}
	execGraph(t, graph)
	assert.Equal(t, name, "b")
}

func TestExecutor_Gen(t *testing.T) {
	graph := NewGraph()
	sum := 0
	sunLock := sync.Mutex{}
	graph.MustNode("a", func(ctx context.Context) ([]*Node, error) {
		nodes := make([]*Node, 0)
		for x := 0; x < 100; x++ {
			x := x
			nodes = append(nodes, MustNode("a", func(ctx context.Context) {
				time.Sleep(time.Millisecond * 10)
				sunLock.Lock()
				sum = sum + x
				sunLock.Unlock()
			}))
		}
		return nodes, nil
	})
	graph.MustNode("b", func(ctx context.Context) {
		assert.Equal(t, sum, 4950)
	})
	if err := NewNodeDependencyBuilder(graph).Node("b").From("a").Build(); err != nil {
		t.Fatal(err)
	}
	exec := NewExecutor(WithLimit(10))
	if err := exec.Execute(context.Background(), graph); err != nil {
		t.Fatal(exec)
	}
}

func TestExecutor_Misc(t *testing.T) {
	graph := NewGraph()

	// 定义 node
	graph.MustNode("a", func(ctx context.Context) string { info(t, "a"); return "b" }) // a返回一个条件分支执行b node
	graph.MustNode("b", func(ctx context.Context) { info(t, "b") })
	graph.MustNode("c", func(ctx context.Context) { info(t, "c") })
	graph.MustNode("d", func(ctx context.Context) []*Node {
		info(t, ContextCurNode(ctx))
		return []*Node{
			MustNode("a", func(ctx context.Context) {
				info(t, ContextCurNode(ctx))
			}),
			MustNode("b", func(ctx context.Context) {
				info(t, ContextCurNode(ctx))
			}),
		}
	})
	graph.MustNode("e", func(ctx context.Context) {
		info(t, "e")
	})

	// 申明 node 关系
	err := NewNodeDependencyBuilder(graph).
		Node("a").To("b", "c").     // a -> b && a -> c
		Node("d").From("b", "c").   // b -> d && c -> d
		Node("e").From("d").Build() // d -> e
	if err != nil {
		t.Fatal(err)
	}

	for x := 0; x < 10; x++ {
		trace := NewTrace()
		// 设置并发度 和 trace
		exec := NewExecutor(WithLimit(x), WithTrace(trace))
		// 执行 node
		if err := exec.Execute(context.Background(), graph); err != nil {
			t.Fatal(err)
		}
		t.Log(trace.Graphviz())
	}
}
