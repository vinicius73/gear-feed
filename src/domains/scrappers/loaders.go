package scrappers

import (
	"gfeed/domains/news"
	"gfeed/domains/news/storage"
	"sync"
	"time"
)

func runWithChannels(wg *sync.WaitGroup, ch chan news.Entry) {
	loaders, err := loadDefinitions()

	if err != nil {
		logger.
			Fatal().
			Err(err).
			Msg("Fail to load definitions")

		return
	}

	wg.Add(len(loaders))

	for _, loader := range loaders {
		go loadIntoChan(wg, ch, loader.FindEnties)
	}

	time.Sleep(time.Second * 1)

	wg.Wait()

	close(ch)
}

func loadIntoChan(wg *sync.WaitGroup, ch chan news.Entry, loader loaderFn) {
	entries := loader()

	for _, v := range entries {
		exist, err := storage.IsRecorded(v)

		if err != nil {
			logger.
				Error().
				Err(err).
				Str("entry", v.Key()).
				Msg("Fail to check record")
		}

		if exist {
			logger.
				Warn().
				Str("scrapper", v.Type).
				Str("title", v.Title).
				Msg("That entry already exist")
		} else {
			ch <- v
		}
	}

	wg.Done()
}
