package main

import (
	"sync"
	"time"
)

type Task struct {
	closed chan struct{}
	wg     sync.WaitGroup
	ticker *time.Ticker
}

func (t *Task) Stop() {
	close(t.closed)
	t.wg.Wait()
}

func (t *Task) Run() {
	for {
		select {
		case <-t.closed:
			return
		case <-t.ticker.C:
			handleChecks()
		}
	}
}
