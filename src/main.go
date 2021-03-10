package main

import (
	"fmt"
	"gfeed/scrappers"
)

func main() {
	entries := scrappers.NewsEntries()

	fmt.Println(entries)
}
