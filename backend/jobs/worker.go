package jobs

import (
	"context"
    "fmt"
	"sync"

	"github.com/charmbracelet/log"
)

func StartWorkerPool(ctx context.Context, jobs <-chan *Job, count int, logger *log.Logger) {
	var wg sync.WaitGroup
	for i := range count {
		wg.Add(1)
        l := logger.WithPrefix(fmt.Sprintf("Worker (%d)", i))

		go func(id int, logger *log.Logger) {
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
        }(i, l)
	}
	wg.Wait()
}
