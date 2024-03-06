package command

import (
	"os"
	"runtime/debug"
	"sync"

	"github.com/anthony-dong/golang/pkg/logs"
)

var deferTask []func()
var deferTaskLock sync.Mutex
var closeOnce sync.Once

func Defer(task func()) {
	if task == nil {
		panic("AddDeferTask: defer task is nil")
	}
	deferTaskLock.Lock()
	deferTask = append(deferTask, task)
	deferTaskLock.Unlock()
}

func Close() {
	closeOnce.Do(func() {
		if len(deferTask) > 0 {
			logs.Builder().Info().String("[defer]").KV("pid", os.Getpid()).String("start run the defer tasks").Emit(nil)
		}
		for index := len(deferTask) - 1; index >= 0; index-- {
			func() {
				defer func() {
					if r := recover(); r != nil {
						logs.Error("run task panic: %v, stack: %s", r, debug.Stack())
					}
				}()
				deferTask[index]()
			}()
		}
	})
}
