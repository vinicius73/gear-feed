package sources

import "github.com/vinicius73/gamer-feed/pkg/scraper"

type Collection []scraper.SourceDefinition

func (c Collection) OnlyStorieSuported() Collection {
	var coll Collection

	for _, source := range c {
		if source.SupportStories {
			coll = append(coll, source)
		}
	}

	return coll
}

func (c Collection) Names() []string {
	names := []string{}

	for _, source := range c {
		names = append(names, source.Name)
	}

	return names
}
