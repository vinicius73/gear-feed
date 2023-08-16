package actions

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/vinicius73/gamer-feed/pkg/stories/stages"
)

type BuildStoryOptions struct {
	URL      string
	Output   string
	Template string
}

func Story(ctx context.Context, opt BuildStoryOptions) error {
	entry, err := stages.Fetch(ctx, opt.URL)
	if err != nil {
		return err
	}

	out, err := filepath.Abs(opt.Output)
	if err != nil {
		return err
	}

	stage, err := stages.BuildStage(stages.BuildStageOptions{
		Source:           entry,
		TargetDir:        out,
		TemplateFilename: opt.Template,
	})
	if err != nil {
		return err
	}

	fmt.Println(stage)

	return nil
}
