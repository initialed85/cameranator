package utils

import (
	"log"

	"github.com/initialed85/glue/pkg/worker"
)

type Work struct {
	WorkFn     func() (interface{}, error)
	CompleteFn func(interface{}, error)
}

type Executor struct {
	blockedWorkers []*worker.BlockedWorker
	queueSize      int
	workQueue      chan Work
}

func NewExecutor(numWorkers int, queueSize int) *Executor {
	e := Executor{
		blockedWorkers: make([]*worker.BlockedWorker, 0),
		queueSize:      queueSize,
	}

	for i := 0; i < numWorkers; i++ {
		e.blockedWorkers = append(e.blockedWorkers, worker.NewBlockedWorker(
			func() {},
			e.work,
			func() {},
		))
	}

	return &e
}

func (e *Executor) work() {
	work, ok := <-e.workQueue

	if !ok {
		return
	}

	result, err := work.WorkFn()
	if err != nil {
		log.Printf("warning: Executor.work; WorkFn() %#+v caused %v", work, err)
	}

	work.CompleteFn(result, err)
}

func (e *Executor) Submit(workFn func() (interface{}, error), completeFn func(interface{}, error)) {
	work := Work{
		WorkFn:     workFn,
		CompleteFn: completeFn,
	}

	e.workQueue <- work
}

func (e *Executor) Start() {
	e.workQueue = make(chan Work, e.queueSize)

	for _, blockedWorker := range e.blockedWorkers {
		blockedWorker.Start()
	}

	log.Printf("Executor.Start; started.")
}

func (e *Executor) Stop() {
	for _, blockedWorker := range e.blockedWorkers {
		blockedWorker.Stop()
	}

	close(e.workQueue)

	log.Printf("Executor.Stop; stopped.")
}
