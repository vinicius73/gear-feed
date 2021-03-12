package scrappers

import (
	log "gfeed/utils/logger"
)

var logger log.Logger

func init() {
	logger = log.New("scrapper")
}
