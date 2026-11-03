package observer

import (
	"sync"
)

type Pool struct {
	jobs    chan func()
	waiting sync.WaitGroup
}

func NewWorkerPool(workers int) *Pool {
	pool := &Pool{jobs: make(chan func())}
	pool.waiting.Add(workers)
	for i := 0; i < workers; i++ {
		go pool.work()
	}
	return pool
}

func (p *Pool) work() {
	defer p.waiting.Done()
	for job := range p.jobs {
		job()
	}
}

func (p *Pool) Submit(job func()) {
	p.jobs <- job
}

func (p *Pool) Shutdown() {
	close(p.jobs)
	p.waiting.Wait()
}
