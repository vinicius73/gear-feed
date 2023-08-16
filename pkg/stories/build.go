package stories

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/vinicius73/gamer-feed/pkg/stories/fetcher"
	"github.com/vinicius73/gamer-feed/pkg/stories/filetemplate"
	"github.com/vinicius73/gamer-feed/pkg/stories/stages"
)

type BuildStorieOptions struct {
	TemplateFilename string
	SourceURL        string
	TargetDir        string
}

func (bo BuildStorieOptions) Template(source fetcher.Result) (filetemplate.Template, error) {
	tpl, err := filetemplate.New(filetemplate.Options{
		Source:   source,
		BaseDir:  bo.TargetDir,
		Template: bo.TemplateFilename,
	})
	if err != nil {
		return filetemplate.Template{}, err
	}

	return tpl, nil
}

func BuildStory(ctx context.Context, opt BuildStorieOptions) (Storie, error) {
	logger := zerolog.Ctx(ctx).With().Str("component", "stories").Logger()
	ctx = logger.WithContext(ctx)

	entry, err := fetcher.Fetch(ctx, fetcher.Options{
		SourceURL:     opt.SourceURL,
		DefaultWidth:  stages.DefaultWidth,
		DefaultHeight: stages.DefaultHeight,
	})
	if err != nil {
		return Storie{}, err
	}

	logger.Info().Str("title", entry.Title).Msg("entry data loades")

	tpl, err := opt.Template(entry)
	if err != nil {
		return Storie{}, err
	}

	stage, err := stages.BuildStage(ctx, stages.BuildStageOptions{
		Source:   entry,
		Template: tpl,
	})
	if err != nil {
		return Storie{}, err
	}

	logger.Info().Str("title", entry.Title).Msg("stage builded")

	videoFile, err := tpl.Render("video.mp4")
	if err != nil {
		return Storie{}, err
	}

	videoFile, err = stages.BuildVideo(ctx, stages.BuildVideoOptions{
		Stage:  stage,
		Target: videoFile,
	})
	if err != nil {
		return Storie{}, err
	}

	logger.Info().Str("title", entry.Title).Msg("video builded")

	return Storie{
		Stage: stage,
		Video: videoFile,
	}, nil
}
