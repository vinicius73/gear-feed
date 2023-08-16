package actions

import (
	"context"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
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

	tmp, err := os.MkdirTemp(os.TempDir(), "gamer-feed-*")
	if err != nil {
		return err
	}

	story, err := stories.BuildStory(ctx, stories.BuildStorieOptions{
		SourceURL:        opt.URL,
		TargetDir:        tmp,
		TemplateFilename: "{{.date}}-{{.site}}-{{.hash}}--{{.filename}}",
	})
	if err != nil {
		return err
	}

	defer story.RemoveStage()

	if err = story.MoveVideo(out); err != nil {
		return err
	}

	logger := zerolog.Ctx(ctx)

	logger.Info().
		Str("hash", story.Hash).
		Str("output", out).Msg("video story was created")

	return nil
}
