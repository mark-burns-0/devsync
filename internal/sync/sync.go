package sync

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Syncer struct {
	dirs []string
}

func New(dirs []string) *Syncer {
	return &Syncer{
		dirs: dirs,
	}
}

func (sy *Syncer) Sync() error {
	process := func(path string) error {
		fmt.Printf("Processing: %s\n", path)
		time.Sleep(1 * time.Second)
		return nil
	}
	queueSize := max(len(sy.dirs)/5, 10)

	wp := newPool(5, queueSize, process)
	wp.Start()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, val := range sy.dirs {
			wp.Add(val)
		}
		wp.CloseJobs()
	}()
	wg.Wait()
	wp.Stop()

	return nil
}

type WorkerPool struct {
	workers   int
	jobs      chan string
	wg        sync.WaitGroup
	stop      chan struct{}
	stopOnce  sync.Once
	closeOnce sync.Once
	running   atomic.Bool
	process   func(string) error
}

func newPool(
	workers int,
	queueSize int,
	process func(string) error,
) *WorkerPool {
	return &WorkerPool{
		workers: workers,
		jobs:    make(chan string, queueSize),
		stop:    make(chan struct{}),
		process: process,
	}
}

func (wp *WorkerPool) Add(dir string) {
	select {
	case wp.jobs <- dir:
	case <-wp.stop:
		return
	}
}

func (wp *WorkerPool) Start() {
	if wp.running.Swap(true) {
		return
	}

	for i := range wp.workers {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	for job := range wp.jobs {
		if err := wp.process(job); err != nil {
			fmt.Println(id, err)
		}
	}
}

func (wp *WorkerPool) Stop() {
	wp.stopOnce.Do(func() {
		wp.running.Swap(false)
		close(wp.stop)
		wp.wg.Wait()
	})
}

func (wp *WorkerPool) CloseJobs() {
	wp.closeOnce.Do(func() {
		close(wp.jobs)
	})
}
