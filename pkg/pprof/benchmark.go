package pprof

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func ParallelBenchmark(name string, thread int, duration time.Duration, execute func(count int)) {
	wg := sync.WaitGroup{}
	wg.Add(thread)
	totalCount := uint64(0)
	totalSpend := uint64(0)
	for i := 0; i < thread; i++ {
		go func() {
			defer wg.Done()
			spend, count := Benchmark(duration, execute)
			atomic.AddUint64(&totalSpend, uint64(spend))
			atomic.AddUint64(&totalCount, uint64(count))
		}()
	}
	wg.Wait()
	fmt.Printf("name=%s thread=%d duration=%s total=%d avg=%s\n", name, thread, duration, totalCount, Avg(time.Duration(totalSpend), int(totalCount)))
}

func Avg(spend time.Duration, count int) string {
	avg := float64(spend) / float64(count)
	if avg > 100 {
		return time.Duration(avg).String()
	}
	return fmt.Sprintf("%.4fns", avg)
}

func Benchmark(duration time.Duration, bench func(count int)) (time.Duration, int) {
	const maxTotalCount = 1000000000 // 10E
	count := 1
	totalSpend := time.Duration(0)
	totalCount := 0
	for {
		start := time.Now()
		bench(count)
		spend := time.Since(start)

		totalSpend = totalSpend + spend
		totalCount = totalCount + count

		if totalCount >= maxTotalCount {
			break
		}
		subSpend := duration - totalSpend
		if subSpend <= 0 {
			break
		}
		count = totalCount*10 - totalCount
		if subCount := int(float64(subSpend) / (float64(totalSpend) / float64(totalCount))); count > subCount {
			count = subCount
		}
	}
	return totalSpend, totalCount
}
