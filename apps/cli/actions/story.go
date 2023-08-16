package actions

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vinicius73/gamer-feed/pkg/stories"
)

type BuildStoryOptions struct {
	URL    string
	Output string
}

func VideoStory(ctx context.Context, opt BuildStoryOptions) error {
	out, err := filepath.Abs(opt.Output)
	if err != nil {
		return err
	}

	tmp, err := os.MkdirTemp(os.TempDir(), "gamer-feed")
	if err != nil {
		return err
	}

	story, err := stories.BuildStory(ctx, stories.BuildStorieOptions{
		SourceURL:        opt.URL,
		TargetDir:        tmp,
		TemplateFilename: "{{.date}}-{{.site}}--{{.filename}}",
	})
	if err != nil {
		return err
	}

	if err = story.MoveVideo(out); err != nil {
		return err
	}

	fmt.Println(story.Video)

	return story.RemoveStage()
}
