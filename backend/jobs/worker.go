package jobs

import (
	"context"
	"sync"
)

func StartWorkerPool(ctx context.Context, jobs <-chan *Job, count int) {
	var wg sync.WaitGroup
	for i := range count {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job := <-jobs:
					job.Run()
					// save result
				}
			}
		}(i)
	}
	wg.Wait()
}
