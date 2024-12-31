package pprof

import (
	"testing"
	"time"
)

func TestBenchmark(t *testing.T) {
	spend, count := Benchmark(time.Second, func(count int) {
		for i := 0; i < count; i++ {
			num := add(1, 2)
			if num != 3 {
				t.Fatalf("add(1,2)=%d\n", num)
			}
		}
	})
	// 1000000000
	// 4023195055
	t.Log(spend, count, Avg(spend, count))
}

func BenchmarkName(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		pb.Next()
	})

}

func add(a, b int) int {
	return a + b
}

func BenchmarkTimeSince(b *testing.B) {
	now := time.Now()
	for i := 0; i < b.N; i++ {
		time.Since(now)
	}
}
