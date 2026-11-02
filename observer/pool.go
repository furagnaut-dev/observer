package observer

import (
	"sync"
)

type Channel struct {
	jobs chan func()
	waiting sync.WaitGroup
}

func NewWorkerPool(workers int) *Channel {
	pool := &Channel{jobs: make(chan func())}
	pool.waiting.Add(workers)
	for i := 0; i < workers; i++ {
		go pool.waiting.Wait()
		defer pool.waiting.Done()
	}
	return pool
}