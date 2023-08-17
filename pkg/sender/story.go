package sender

import (
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/stories"
)

type Story[T model.IEntry] struct {
	Story stories.Story
	Entry T
}
