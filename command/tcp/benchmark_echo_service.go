package tcp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"net"
	"sync/atomic"
	"time"

	"github.com/anthony-dong/golang/pkg/utils"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/anthony-dong/golang/pkg/logs"
)

func NewBenchmarkEchoServiceCommand() (*cobra.Command, error) {
	addr := ""
	maxConn := 0
	concurrent := 0
	payloadSize := 0
	interval := time.Duration(0)
	runTime := time.Duration(0)
	command := &cobra.Command{
		Use:   "benchmark_echo_service",
		Short: "benchmark echo service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return benchmark(cmd.Context(), addr, maxConn, concurrent, newBuffer(payloadSize), runTime, interval)
		},
	}
	command.Flags().StringVar(&addr, "addr", "localhost:8080", "请求地址")
	command.Flags().IntVar(&maxConn, "conn", 10, "最大连接数")
	command.Flags().IntVar(&concurrent, "count", 0, "最大并发请求(0:表示不限制)")
	command.Flags().DurationVar(&runTime, "run", time.Second*10, "最大请求时间")
	command.Flags().IntVar(&payloadSize, "body", 64, "请求包大小")
	command.Flags().DurationVar(&interval, "interval", 0, "请求间隔")
	return command, nil
}

type Limiter struct {
	buf chan bool
}

func newLimiter(size int) *Limiter {
	if size <= 0 {
		return &Limiter{nil}
	}
	return &Limiter{make(chan bool, size)}
}

func (l *Limiter) acquire() {
	if l == nil || l.buf == nil {
		return
	}
	l.buf <- true
}

func (l *Limiter) release() {
	if l == nil || l.buf == nil {
		return
	}
	<-l.buf
}

type Counter struct {
	TotalSpend      uint64
	TotalSucRequest uint64
	TotalErrRequest uint64
}

func benchmark(ctx context.Context, addr string, maxConn int, concurrent int, writeData []byte, runTime, interval time.Duration) error {
	logs.CtxInfo(ctx, "addr=%s, max_conn=%v, concurrent=%v, data_size=%v, run_time=%s, interval=%s", addr, maxConn, concurrent, len(writeData), runTime, interval)
	limiter := newLimiter(concurrent)
	counter := &Counter{}
	realAddr, err := utils.ParseAddr(addr)
	if err != nil {
		return err
	}
	logs.CtxInfo(ctx, "dial addr: %s", realAddr)
	stopTime := time.Now().Add(runTime)
	group := errgroup.Group{}
	start := time.Now()
	for x := 0; x < maxConn; x++ {
		group.Go(func() error {
			conn, err := realAddr.Dial()
			if err != nil {
				return err
			}
			defer conn.Close()
			count := 0
			for {
				if time.Now().After(stopTime) {
					return nil
				}
				count = count + 1
				wrapperRequest(limiter, counter, interval, conn, writeData, count == 1)
			}
		})
	}
	defer func() {
		if counter.TotalSucRequest > 0 && counter.TotalSpend > 0 {
			logs.Info("latency avg(req): %s\n", time.Duration(counter.TotalSpend/counter.TotalSucRequest))
		}
		if counter.TotalSucRequest > 0 || counter.TotalErrRequest > 0 {
			logs.Info("throughput avg(s): %d\n", (counter.TotalErrRequest+counter.TotalSucRequest)/uint64(time.Since(start).Seconds()))
		}
		logs.Info("total success request: %d\n", counter.TotalSucRequest)
		logs.Info("total error request: %d\n", counter.TotalErrRequest)
	}()
	if err := group.Wait(); err != nil {
		return err
	}
	return nil
}

func newBuffer(size int) []byte {
	output := make([]byte, size)
	for x := 0; x < size; x++ {
		output[x] = byte(rand.Int31()%127) + 1
	}
	return output
}

func wrapperRequest(limiter *Limiter, counter *Counter, interval time.Duration, conn net.Conn, writeData []byte, skipRecord bool) {
	limiter.acquire()
	defer limiter.release()

	if interval > 0 {
		time.Sleep(interval)
	}
	start := time.Now()
	if err := request(conn, writeData); err != nil {
		atomic.AddUint64(&counter.TotalErrRequest, 1)
		logs.Error("request conn [%s] find err: %v", conn.RemoteAddr(), err)
		return
	}
	if skipRecord {
		return
	}
	atomic.AddUint64(&counter.TotalSpend, uint64(time.Now().Sub(start)))
	atomic.AddUint64(&counter.TotalSucRequest, 1)
}

func request(conn net.Conn, writeData []byte) error {
	_, err := conn.Write(writeData)
	if err != nil {
		return err
	}
	readData := make([]byte, len(writeData))
	if _, err := io.ReadFull(conn, readData); err != nil {
		return err
	}
	if !bytes.Equal(writeData, readData) {
		return fmt.Errorf(`not equal`)
	}
	return nil
}
