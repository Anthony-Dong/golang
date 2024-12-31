package pprof

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
)

// InitPProf
// go InitPProf()
func InitPProf() {
	err := http.ListenAndServe(":12345", http.DefaultServeMux)
	if err != nil {
		panic(err)
	}
}

func StartCPUProfile(fileName string) (stop func()) {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		if err := f.Close(); err != nil {
			panic(err)
		}
		panic(err)
	}
	return func() {
		pprof.StopCPUProfile()
		if err := f.Close(); err != nil {
			panic(err)
		}
	}
}

func StartMemProfile(fileName string) (stop func()) {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	return func() {
		defer func() {
			if err := f.Close(); err != nil {
				panic(err)
			}
		}()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			panic(err)
		}
	}
}

func StartTraceProfile(fileName string) (stop func()) {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	if err := trace.Start(f); err != nil {
		panic(err)
	}
	return func() {
		defer func() {
			trace.Stop()
			if err := f.Close(); err != nil {
				panic(err)
			}
		}()
	}
}
