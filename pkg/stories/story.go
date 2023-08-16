package stories

import (
	"os"

	"github.com/vinicius73/gamer-feed/pkg/stories/stages"
)

type Story struct {
	Stage stages.Stage
	Video string
	Hash  string
}

func (s Story) MoveVideo(target string) error {
	return os.Rename(s.Video, target)
}

func (s Story) RemoveStage() error {
	return s.Stage.RemoveAll()
}
