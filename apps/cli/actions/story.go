package actions

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/vinicius73/gamer-feed/pkg/stories/covers"
)

type BuildStoryOptions struct {
	URL      string
	Output   string
	Template string
}

func Story(ctx context.Context, opt BuildStoryOptions) error {
	entry, err := covers.Fetch(ctx, opt.URL)
	if err != nil {
		return err
	}

	out, err := filepath.Abs(opt.Output)
	if err != nil {
		return err
	}

	res, err := covers.NewBuilder(entry).Build(covers.BuildOptions{
		TargetDir:        out,
		TemplateFilename: opt.Template,
	})
	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}
