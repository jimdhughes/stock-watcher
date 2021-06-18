package util

import (
	"sync"
	"time"
)

type Task struct {
	Closed chan struct{}
	Wg     sync.WaitGroup
	Ticker *time.Ticker
}

func (t *Task) Stop() {
	close(t.Closed)
	t.Wg.Wait()
}

func (t *Task) Run(fn func()) {
	for {
		select {
		case <-t.Closed:
			return
		case <-t.Ticker.C:
			fn()
		}
	}
}
