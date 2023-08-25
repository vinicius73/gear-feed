package actions

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/vinicius73/gamer-feed/pkg/configurations"
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/sources"
	"github.com/vinicius73/gamer-feed/pkg/stories"
	"github.com/vinicius73/gamer-feed/pkg/tasks"
)

type BuildStoryOptions struct {
	URL    string
	Output string
	Footer stories.Footer
}

type SendStoriesOptions struct {
	Sources sources.LoadOptions
	Footer  stories.Footer
	Period  time.Duration
	To      int64
	Limit   int
}

func VideoStory(ctx context.Context, opt BuildStoryOptions) error {
	out, err := filepath.Abs(opt.Output)
	if err != nil {
		return err
	}

	tmp, err := os.MkdirTemp(os.TempDir(), "gfeed-*")
	if err != nil {
		return err
	}

	story, err := stories.BuildStory(ctx, stories.BuildStorieOptions{
		SourceURL:        opt.URL,
		TargetDir:        tmp,
		Footer:           opt.Footer,
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

func SendStories(ctx context.Context, opt SendStoriesOptions) error {
	config := configurations.Ctx(ctx)

	store, db, err := buildDB[model.Entry](ctx, config)
	if err != nil {
		return err
	}

	defer db.Close()

	botSender, err := buildSender(SenderOptions{
		Chats:    []int64{opt.To},
		Storage:  store,
		Telegram: config.Telegram,
	})
	if err != nil {
		return err
	}

	return tasks.SendLastStories[model.Entry]{
		Limit:    opt.Limit,
		Sources:  opt.Sources,
		Interval: opt.Period,
		Footer:   opt.Footer,
	}.
		Run(ctx, tasks.TaskRunOptions[model.Entry]{
			Storage: store,
			Sender:  botSender,
		})
}
