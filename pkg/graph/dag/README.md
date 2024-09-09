# Dag 引擎


```go
func TestExecutor_Misc(t *testing.T) {
	graph := NewGraph()

	// 定义 node
	graph.MustNode("a", func(ctx context.Context) string { info(t, "a"); return "b" })
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
	err := graph.NewDepBuilder().
		Node("a").Before("b", "c").  // a -> b && a -> c
		Node("d").After("b", "c").   // b -> d && c -> d
		Node("e").After("d").Build() // d -> e
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
```