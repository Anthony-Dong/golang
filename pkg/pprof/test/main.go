package main

import (
	"sync"
	"time"

	"github.com/anthony-dong/golang/pkg/pprof"
)

func main() {
	// 记录 cup pprof
	//stop := pprof.StartCPUProfile("cpu.out")
	//defer stop()

	// 并发测试 sync map的性能
	mm := sync.Map{}
	pprof.ParallelBenchmark("test1", 64, time.Second, func(count int) {
		for i := 0; i < count; i++ {
			mm.Store(i%10000, 1)
		}
	})
	// name=test1 thread=32 duration=1s total=6708009 avg=4.772µs
	// name=test1 thread=64 duration=1s total=6883456 avg=9.3µs
}
