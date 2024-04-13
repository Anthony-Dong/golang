package tcp

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net"
	"sync/atomic"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/anthony-dong/golang/pkg/logs"
)

func NewBenchmarkEchoServiceCommand() (*cobra.Command, error) {
	addr := ""
	maxConn := 0
	maxReq := 0
	payloadSize := 0
	interval := time.Duration(0)
	runTime := time.Duration(0)
	command := &cobra.Command{
		Use:   "benchmark_echo_service",
		Short: "benchmark echo service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return benchmark(addr, maxConn, maxReq, newBuffer(payloadSize), runTime, interval)
		},
	}
	command.Flags().StringVar(&addr, "addr", "localhost:8080", "请求地址")
	command.Flags().IntVar(&maxConn, "conn", 10, "最大连接数")
	command.Flags().IntVar(&maxReq, "count", 1000, "最大并发请求")
	command.Flags().DurationVar(&runTime, "run", time.Second*10, "最大请求时间")
	command.Flags().IntVar(&payloadSize, "body", 64, "请求包大小")
	command.Flags().DurationVar(&interval, "interval", 0, "请求间隔")
	return command, nil
}

func benchmark(addr string, maxConn int, maxReq int, writeData []byte, runTime, interval time.Duration) error {
	logs.Info("addr=%s, max_conn=%v, max_req=%v, data_size=%v, run_time=%s, interval=%s", addr, maxConn, maxReq, len(writeData), runTime, interval)
	buffers := make(chan bool, maxReq)
	totalSpend := uint64(0)
	totalRequest := uint64(0)
	totalErrRequest := uint64(0)
	stopTime := time.Now().Add(runTime)
	group := errgroup.Group{}
	for x := 0; x < maxConn; x++ {
		group.Go(func() error {
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				return err
			}
			defer conn.Close()
			count := 0
			for {
				if time.Now().After(stopTime) {
					return nil
				}
				buffers <- true

				if interval > 0 {
					time.Sleep(interval)
				}
				start := time.Now()
				if err := request(conn, writeData); err != nil {
					atomic.AddUint64(&totalErrRequest, 1)
					logs.Error("request conn [%s] find err: %v", conn.RemoteAddr(), err)

					<-buffers
					continue
				}
				// 第一次连接不算
				count++
				if count == 1 {
					continue
				}
				atomic.AddUint64(&totalSpend, uint64(time.Now().Sub(start)))
				atomic.AddUint64(&totalRequest, 1)

				<-buffers
			}
		})
	}
	defer func() {
		if totalRequest > 0 && totalSpend > 0 {
			logs.Info("latency avg(req): %s\n", time.Duration(totalSpend/totalRequest))
		}
		if totalRequest > 0 {
			logs.Info("throughput avg(s): %d\n", totalRequest/uint64(runTime/time.Second))
		}
		logs.Info("total success request: %d\n", totalRequest)
		logs.Info("total error request: %d\n", totalErrRequest)
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
