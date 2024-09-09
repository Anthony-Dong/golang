package dag

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"strconv"
)

type Executor interface {
	Execute(ctx context.Context, graph *Graph) error
}

// Executor todo set concurrent limit
type defaultExecutor struct {
	trace Trace // todo OpenTrace

	results      map[string]*executeResult
	completeChan chan *executeResult

	nodeFroms map[string][]string // node: from nodes, 解决动态生成node，保存 pre_node.to_node 的 from nodes

	execLimiter chan bool
}

// WithTrace set trace
func WithTrace(trace Trace) func(executor *defaultExecutor) {
	return func(executor *defaultExecutor) {
		executor.trace = trace
	}
}

// WithLimit set concurrent limit
func WithLimit(limit int) func(executor *defaultExecutor) {
	return func(executor *defaultExecutor) {
		if limit <= 0 {
			executor.execLimiter = nil
			return
		}
		executor.execLimiter = make(chan bool, limit)
	}
}

func NewExecutor(ops ...func(executor *defaultExecutor)) Executor {
	exec := &defaultExecutor{}
	for _, op := range ops {
		if op == nil {
			continue
		}
		op(exec)
	}
	return exec
}

type executeResult struct {
	Node *Node
	Err  error
	Cond string
	Gen  []*Node
}

func (e *executeResult) Error() string {
	return fmt.Sprintf(`exec node [%s] find err: %s`, e.Node.Name, e.Err)
}

func (e *defaultExecutor) init(graph *Graph) error {
	e.results = make(map[string]*executeResult, len(graph.nodes))
	e.completeChan = make(chan *executeResult, len(graph.nodes))
	e.nodeFroms = make(map[string][]string)
	return nil
}

func (e *defaultExecutor) Execute(ctx context.Context, graph *Graph) error {
	if err := e.init(graph); err != nil {
		return err
	}
	if ctx == nil {
		ctx = context.Background()
	}

	// exec begin node
	e.exec(ctx, graph.begin)

	for {
		select {
		case result := <-e.completeChan:
			e.results[result.Node.Name] = result

			ctx := setCtxFromNode(ctx, result.Node.Name)

			// run sub node
			var err error
			switch result.Node.GetRunType() {
			case NodeRunTypeNormal:
				err = e.succeedRunNode(ctx, result)
			case NodeRunTypeCond:
				err = e.succeedRunCondNode(ctx, result)
			case NodeRunTypeGen:
				err = e.succeedRunGenNode(ctx, result)
			default:
				err = e.succeedRunNode(ctx, result)
			}
			if err != nil {
				return err
			}
			if e.isDone(graph) {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (e *defaultExecutor) mergeNodeFroms(node *Node) []string {
	result := make([]string, 0, len(node.From)+len(e.nodeFroms[node.Name]))
	for _, elem := range node.From {
		result = append(result, elem.Name)
	}
	for _, elem := range e.nodeFroms[node.Name] {
		result = append(result, elem)
	}
	return result
}

func (e *defaultExecutor) canRunNode(node *Node) (bool, error) {
	count := 0
	froms := e.mergeNodeFroms(node)
	for _, from := range froms {
		result := e.results[from]
		if result == nil {
			continue
		}
		count++
		if result.Err != nil {
			if node.Downgrade[result.Node.Name] {
				continue
			}
			return false, result
		}
	}
	return count == len(froms), nil
}

func (e *defaultExecutor) isDone(graph *Graph) bool {
	return e.results[graph.end.Name] != nil
}

func (e *defaultExecutor) succeedRunNode(ctx context.Context, result *executeResult) error {
	for _, to := range result.Node.To {
		can, err := e.canRunNode(to)
		if err != nil {
			return err
		}
		if can {
			e.exec(ctx, to)
		}
	}
	return nil
}

func (e *defaultExecutor) succeedRunGenNode(ctx context.Context, result *executeResult) error {
	mm := make(map[string]int, len(result.Gen))
	for _, node := range result.Gen {
		// 防止sub node name 重复了
		mm[node.Name]++
		if count := mm[node.Name]; count > 1 {
			node.Name = node.Name + "_" + strconv.Itoa(count)
		}

		// join 父亲节点的name
		node.Name = result.Node.Name + "." + node.Name
		node.To = result.Node.To
		for _, to := range node.To {
			e.nodeFroms[to.Name] = append(e.nodeFroms[to.Name], node.Name)
		}
		e.exec(ctx, node)
	}
	return nil
}

func (e *defaultExecutor) succeedRunCondNode(ctx context.Context, result *executeResult) error {
	found := false
	for _, to := range result.Node.To {
		if to.Name != result.Cond {
			e.results[to.Name] = &executeResult{Node: to} // 由于当前节点永远也执行不到，所以这里直接success
			continue
		}
		found = true
		can, err := e.canRunNode(to)
		if err != nil {
			return err
		}
		if can {
			e.exec(ctx, to)
		}
	}
	if !found {
		return fmt.Errorf(`not found cond node: %s`, result.Cond)
	}
	return nil
}

func (e *defaultExecutor) exec(ctx context.Context, node *Node) {
	if e.execLimiter != nil {
		e.execLimiter <- true
	}
	go func() {
		var (
			cond string
			err  error
			gen  []*Node
		)
		ctx := setCtxCurNode(ctx, node.Name)
		defer func() {
			// <-e.execLimiter 必须要在 e.completeChan<- 之前
			// 不然当 execLimiter buffer满了，当e.completeChan <- executeResult时主线程在写阻塞e.execLimiter和读阻塞e.completeChan，死锁了!
			if e.execLimiter != nil {
				<-e.execLimiter
			}
			if r := recover(); r != nil {
				log.Printf("panic: %v\n%s", r, debug.Stack())
				err = fmt.Errorf(`panic: %v`, r)
			}
			if e.trace != nil {
				e.trace.Report(contextFromNode(ctx), node.Name, nil)
			}
			e.completeChan <- &executeResult{
				Node: node,
				Err:  err,
				Cond: cond,
				Gen:  gen,
			}
		}()
		switch node.GetRunType() {
		case NodeRunTypeNormal:
			err = node.Run(ctx)
		case NodeRunTypeCond:
			cond, err = node.RunCond(ctx)
		case NodeRunTypeGen:
			gen, err = node.RunGen(ctx)
		default:
			//
		}
	}()
}

func contextFromNode(ctx context.Context) string {
	v, _ := ctx.Value("ContextFromNode").(string)
	return v
}

func ContextCurNode(ctx context.Context) string {
	v, _ := ctx.Value("ContextCurNode").(string)
	return v
}

func setCtxFromNode(ctx context.Context, from string) context.Context {
	return context.WithValue(ctx, "ContextFromNode", from)
}

func setCtxCurNode(ctx context.Context, cur string) context.Context {
	return context.WithValue(ctx, "ContextCurNode", cur)
}
