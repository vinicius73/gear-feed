package sourcelist

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/vinicius73/gear-feed/pkg/scraper"
)

var _ list.DefaultItem = (*SourceItem)(nil)

type SourceItem struct {
	scraper.SourceDefinition
}

func (i SourceItem) Title() string       { return i.Name }
func (i SourceItem) Description() string { return i.BaseURL }
func (i SourceItem) FilterValue() string { return i.Name }
