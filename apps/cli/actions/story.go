package actions

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/vinicius73/gamer-feed/pkg/stories"
)

type BuildStoryOptions struct {
	URL      string
	Output   string
	Template string
}

func Story(ctx context.Context, opt BuildStoryOptions) error {
	out, err := filepath.Abs(opt.Output)
	if err != nil {
		return err
	}

	stage, err := stories.BuildStory(ctx, stories.BuildStorieOptions{
		SourceURL:        opt.URL,
		TargetDir:        out,
		TemplateFilename: opt.Template,
	})
	if err != nil {
		return err
	}

	fmt.Println(stage)

	return nil
}
