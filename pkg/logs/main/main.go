package main

import "github.com/anthony-dong/golang/pkg/logs"

func main() {
	logs.SetLevel(logs.LevelDebug)
	ctx := logs.CtxWithLogID(nil, "1111")
	logs.CtxDebug(ctx, "hello %v", "world")
	logs.CtxInfo(ctx, "hello %v", "world")
	logs.CtxWarn(ctx, "hello %v", "world")
	logs.CtxError(ctx, "hello %v", "world")
}
