package scrappers

import (
	"gfeed/news"
	"sync"
)

// NewsEntries Load last news Entries
func NewsEntries() (entries []news.Entry) {
	var wg sync.WaitGroup

	ch := make(chan news.Entry, 100)

	logger.Info().Msg("Starting scrappers...")

	runWithChannels(&wg, ch)

	logger.Info().Msg("Scrappers are finish...")

	for entry := range ch {
		logger.
			Info().
			Str("type", entry.Type).
			Msgf("Entry: %s", entry.Link)

		entries = append(entries, entry)
	}

	logger.Info().Msg("Done.")

	return entries
}
