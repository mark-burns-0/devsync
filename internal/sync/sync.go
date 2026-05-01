package sync

import (
	"sync"
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
	queueSize := max(len(sy.dirs)/5, 10)

	wp := newPool(5, queueSize, syncGitRepository)
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
