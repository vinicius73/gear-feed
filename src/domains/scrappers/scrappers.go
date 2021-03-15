package scrappers

import (
	"gfeed/domains/news"
	"sync"
	"time"
)

// NewsEntries Load last news Entries
func NewsEntries() (entries []news.Entry) {
	var wg sync.WaitGroup

	ch := make(chan news.Entry, 100)

	startTime := time.Now()

	logger.Info().Msg("Starting scrappers...")

	runWithChannels(&wg, ch)

	logger.
		Info().
		Msg("Scrappers are done.")

	for entry := range ch {
		logger.
			Info().
			Str("type", entry.Type).
			Msgf("Entry: %s", entry.Link)

		entries = append(entries, entry)
	}

	logger.
		Info().
		Msgf("Done (%s).", time.Since(startTime).String())

	return entries
}
