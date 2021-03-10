package scrappers

import "gfeed/utils"

var logger utils.Logger

func init() {
	logger = utils.NewLogger("scrapper")
}
