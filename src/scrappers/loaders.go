package scrappers

import (
	"gfeed/news"
	"gfeed/news/data"
	"gfeed/scrappers/loader"
	"sync"
	"time"
)

type loaderFn = func() []news.Entry

var loaders []loaderFn

func init() {
	ign := loader.Definitions{
		BaseURL: "https://br.ign.com",
		Path:    "/",
		Name:    "IGN",
		Limit:   2,
		Attributes: loader.AttributesFinder{
			Wrapper: ".hspotlight .card",
			Title: loader.PathFinder{
				Path: ".caption > a",
			},
			Link: loader.PathFinder{
				Path:      ".caption > a",
				Attribute: "href",
			},
			Image: loader.PathFinder{
				Path:      "img",
				Attribute: "src",
			},
		},
	}

	tecnoblog := loader.Definitions{
		BaseURL: "https://tecnoblog.net",
		Path:    "/cat/games-jogos/",
		Name:    "TECNOBLOG",
		Limit:   2,
		Attributes: loader.AttributesFinder{
			Wrapper: "div.posts article.bloco",
			Title: loader.PathFinder{
				Path: ".texts h2 a",
			},
			Link: loader.PathFinder{
				Path:      ".texts a",
				Attribute: "href",
			},
			Image: loader.PathFinder{
				Path:      ".thumb img",
				Attribute: "data-lazy-src",
			},
			Category: loader.PathFinderCategory{
				Alloweds: []string{"news", "especial", "hit kill"},
				PathFinder: loader.PathFinder{
					Path: ".thumb .spread a",
				},
			},
		},
	}

	voxel := loader.Definitions{
		BaseURL: "https://www.tecmundo.com.br",
		Path:    "/voxel",
		Name:    "VOXEL",
		Limit:   2,
		Attributes: loader.AttributesFinder{
			Wrapper: "article.tec--voxel-main__item",
			Title: loader.PathFinder{
				Path: ".tec--voxel-main__item__title a.tec--voxel-main__item__title__link",
			},
			Link: loader.PathFinder{
				Path:      "figure > a",
				Attribute: "href",
			},
			Image: loader.PathFinder{
				Path:      "img.tec--voxel-main__item__thumb__image",
				Attribute: "data-src",
			},
		},
	}

	theenemy := loader.Definitions{
		BaseURL: "https://www.theenemy.com.br",
		Path:    "/news",
		Name:    "THE_ENEMY",
		Limit:   2,
		Attributes: loader.AttributesFinder{
			Wrapper: "div[data-content-type='featured'] .card__wrapper",
			Title: loader.PathFinder{
				Path: ".card__content > a",
			},
			Link: loader.PathFinder{
				Path:      ".card__content > a",
				Attribute: "href",
			},
			Image: loader.PathFinder{
				Path:          "a.card__image__anchor",
				Attribute:     "style",
				ParseStrategy: loader.ParserStrategyStyle,
			},
		},
	}

	loaders = []loaderFn{theenemy.FindEnties, ign.FindEnties, tecnoblog.FindEnties, voxel.FindEnties}
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
				Error().
				Err(err).
				Str("entry", v.Key()).
				Msg("Fail to check record")
		}

		if exist {
			logger.
				Warn().
				Str("entry", v.Key()).
				Msg("That entry already exist")
		} else {
			ch <- v
		}
	}

	wg.Done()
}
