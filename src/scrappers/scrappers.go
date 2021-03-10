package scrappers

import (
	"gfeed/news"
	"sync"
)

// NewsEntries Load last news Entries
func NewsEntries() (entries []news.Entry) {
	var wg sync.WaitGroup

	ch := make(chan news.Entry, 2)

	runOverChanners(&wg, ch)

	for entry := range ch {
		entries = append(entries, entry)
	}

	return entries
}
