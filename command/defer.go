package command

import (
	"runtime/debug"
	"sync"

	"github.com/anthony-dong/golang/pkg/logs"
)

var deferTask []func()
var deferTaskLock sync.Mutex
var closeSignalOnce sync.Once

func AddDeferTask(task func()) {
	if task == nil {
		panic("AddDeferTask: defer task is nil")
	}
	deferTaskLock.Lock()
	deferTask = append(deferTask, task)
	deferTaskLock.Unlock()
}

func CloseDeferTask() {
	closeSignalOnce.Do(func() {
		for index := len(deferTask) - 1; index >= 0; index-- {
			func() {
				defer func() {
					if r := recover(); r != nil {
						logs.Error("CloseDeferTask panic: %v, stack: %s", r, debug.Stack())
					}
				}()
				deferTask[index]()
			}()
		}
	})
}
