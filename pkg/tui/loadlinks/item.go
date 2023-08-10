package loadlinks

import "github.com/vinicius73/gamer-feed/pkg/model"

type Link struct {
	model.Entry
}

func (i Link) Title() string       { return i.Entry.Title }
func (i Link) Description() string { return i.URL }
func (i Link) FilterValue() string { return i.Entry.Title }
