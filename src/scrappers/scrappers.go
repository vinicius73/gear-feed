package scrappers

import (
	"gfeed/news"
	"sync"
)

// NewsEntries Load last news Entries
func NewsEntries() (entries []news.Entry) {
	var wg sync.WaitGroup

	ch := make(chan news.Entry, 100)

	logger.Info("Starting scrappers...")

	runWithChannels(&wg, ch)

	logger.Info("Scrappers are finish...")

	for entry := range ch {
		logger.WithField("type", entry.Type).Info("Entry: " + entry.Link)
		entries = append(entries, entry)
	}

	logger.Info("Done")

	return entries
}
