package scrappers

import (
	"gfeed/news"
	"gfeed/scrappers/theenemy"
	"sync"
)

type loaderFn = func() []news.Entry

var loaders []loaderFn

func init() {
	loaders = []loaderFn{theenemy.Load}
}

func runOverChanners(wg *sync.WaitGroup, ch chan news.Entry) {
	for _, loader := range loaders {
		wg.Add(1)

		go loadIntoChan(wg, ch, loader)
	}

	wg.Wait()

	close(ch)
}

func loadIntoChan(wg *sync.WaitGroup, ch chan news.Entry, loader loaderFn) {
	entries := loader()

	for _, v := range entries {
		ch <- v
	}

	wg.Done()
}
