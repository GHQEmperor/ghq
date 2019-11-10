package watchdogs

import (
	"sync"
	"time"
)

type WatchDog struct {
	key            interface{}
	context        interface{}
	foodChan       chan struct{}
	deadChan       chan struct{}
	deadReturnChan chan struct{}
	destroyChan    chan struct{}
	isClose        bool
	mutex          sync.Mutex

	container *Container
}

func (w *WatchDog) GetKey() interface{} {
	return w.key
}

func (w *WatchDog) feed() {
	w.mutex.Lock()
	if !w.isClose {
		w.foodChan <- struct{}{}
	}
	w.mutex.Unlock()
}

func (w *WatchDog) dead() {
	w.mutex.Lock()
	if !w.isClose {
		w.deadChan <- struct{}{}
	}
	w.mutex.Unlock()
}

func (w *WatchDog) deadReturn() {
	w.mutex.Lock()
	if !w.isClose {
		w.deadReturnChan <- struct{}{}
	}
	w.mutex.Unlock()
}

func (w *WatchDog) destroy() {
	w.mutex.Lock()
	if !w.isClose {
		w.destroyChan <- struct{}{}
	}
	w.mutex.Unlock()
}

func (w *WatchDog) close() {
	w.mutex.Lock()
	if !w.isClose {
		w.isClose = true
		close(w.foodChan)
		close(w.deadChan)
		close(w.deadReturnChan)
		close(w.destroyChan)
	}
	w.mutex.Unlock()
}

/*
	deadFunc, deadReturn
*/
func (w *WatchDog) dogRunning(timeoutFunc, deadFunc, deadReturn func(context interface{}), timeout time.Duration) {
	afterTime := timeout
	after := time.NewTimer(afterTime)
	defer after.Stop()
	for {
		after.Reset(afterTime)
		select {
		case <-w.foodChan:
		case <-after.C:
			timeoutFunc(w.context)
			w.container.removeDog(w.key)
			w.close()
			return
		case <-w.deadChan:
			deadFunc(w.context)
			w.container.removeDog(w.key)
			w.close()
			return
		case <-w.deadReturnChan:
			deadReturn(w.context)
			w.container.removeDog(w.key)
			w.close()
			return
		case <-w.destroyChan:
			w.container.removeDog(w.key)
			w.close()
			return
		}
	}
}
