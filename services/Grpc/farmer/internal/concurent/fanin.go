package concurent

import (
	"context"
	"sync"
)

func multiplex(
	ctx context.Context,
	wg *sync.WaitGroup,
	rcv chan<- any,
	c <-chan any,
) {
	defer wg.Done()
	for v := range c {
		select {
		case <-ctx.Done():
			return // berhenti jika sudah ada sinyal done
		case rcv <- v:
		}
	}
}

func FanIn(
	ctx context.Context,
	channels ...<-chan any,
) <-chan any {
	var wg sync.WaitGroup
	out := make(chan any, len(channels))

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(ctx, &wg, out, c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
