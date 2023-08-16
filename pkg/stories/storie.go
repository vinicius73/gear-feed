package stories

import (
	"os"

	"github.com/vinicius73/gamer-feed/pkg/stories/stages"
)

type Storie struct {
	Stage stages.Stage
	Video string
}

func (s Storie) MoveVideo(target string) error {
	return os.Rename(s.Video, target)
}

func (s Storie) RemoveStage() error {
	return s.Stage.RemoveAll()
}
