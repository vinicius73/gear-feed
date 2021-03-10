package main

import (
	"fmt"
	"gfeed/scrappers"
)

func main() {
	entries := scrappers.NewsEntries()

	for _, v := range entries {
		fmt.Println(v)
	}
}
