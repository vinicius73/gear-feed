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

type Collection []Story

func (s Story) MoveVideo(target string) error {
	return os.Rename(s.Video, target)
}

func (s Story) RemoveStage() error {
	return s.Stage.RemoveAll()
}

func (s Story) RemoveAll() error {
	if err := s.RemoveStage(); err != nil {
		return err
	}

	return os.Remove(s.Video)
}

func (c Collection) RemoveAll() error {
	for _, story := range c {
		if err := story.RemoveAll(); err != nil {
			return err
		}
	}

	return nil
}
