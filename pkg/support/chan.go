package support

import "sync"

func MergeChanners[T any](list ...<-chan T) <-chan T {
	var wg sync.WaitGroup

	out := make(chan T)

	output := func(c <-chan T) {
		defer wg.Done()

		for n := range c {
			out <- n
		}
	}

	wg.Add(len(list))

	for _, c := range list {
		go output(c)
	}

	go func() {
		wg.Wait()

		close(out)
	}()

	return out
}
