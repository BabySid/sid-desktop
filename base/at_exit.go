package base

import (
	"sync"
)

var (
	atExitHandlers   = make([]func(), 0)
	atExitHandleLock sync.RWMutex

	atExitOnce sync.Once
)

func Exit() {
	atExitOnce.Do(func() {
		executeHandlers()
	})
}

func RegisterAtExit(handler func()) {
	atExitHandleLock.Lock()
	defer atExitHandleLock.Unlock()

	atExitHandlers = append(atExitHandlers, handler)
}

func runHandler(handler func()) {
	defer func() {
		recover()
	}()

	handler()
}

func executeHandlers() {
	atExitHandleLock.RLock()
	defer atExitHandleLock.RUnlock()

	size := len(atExitHandlers)
	for i := size - 1; i >= 0; i-- {
		runHandler(atExitHandlers[i])
	}
}
