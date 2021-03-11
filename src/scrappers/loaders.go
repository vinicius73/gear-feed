package scrappers

import (
	"gfeed/news"
	"gfeed/news/data"
	"gfeed/scrappers/tecnoblog"
	"gfeed/scrappers/theenemy"
	"gfeed/scrappers/voxel"
	"sync"
	"time"
)

type loaderFn = func() []news.Entry

var loaders []loaderFn

func init() {
	loaders = []loaderFn{theenemy.Load, voxel.Load, tecnoblog.Load}
}

func runWithChannels(wg *sync.WaitGroup, ch chan news.Entry) {
	wg.Add(len(loaders))

	for _, loader := range loaders {
		go loadIntoChan(wg, ch, loader)
	}

	time.Sleep(time.Second * 1)

	wg.Wait()

	close(ch)
}

func loadIntoChan(wg *sync.WaitGroup, ch chan news.Entry, loader loaderFn) {
	entries := loader()

	for _, v := range entries {
		exist, err := data.IsRecorded(v)

		if err != nil {
			logger.
				WithField("entry", v.Key()).
				Error("data.IsRecorded: " + err.Error())
		}

		if exist {
			logger.
				WithField("entry", v.Key()).
				Warn("That entry already exist")
		} else {
			ch <- v
		}
	}

	wg.Done()
}
