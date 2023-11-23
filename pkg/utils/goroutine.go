package utils

import (
	"runtime/debug"

	"github.com/anthony-dong/golang/pkg/logs"
)

func GoRecoverFunc(foo func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("GoRecoverFunc panic: %v, stack: %s", r, debug.Stack())
			}
		}()
		foo()
	}()
}
