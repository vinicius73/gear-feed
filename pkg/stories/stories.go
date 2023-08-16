package stories

import (
	"context"

	"github.com/vinicius73/gamer-feed/pkg/stories/fetcher"
	"github.com/vinicius73/gamer-feed/pkg/stories/filetemplate"
	"github.com/vinicius73/gamer-feed/pkg/stories/stages"
)

type Storie struct {
	Stage stages.Stage
	Video string
}

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
	entry, err := fetcher.Fetch(ctx, fetcher.Options{
		SourceURL:     opt.SourceURL,
		DefaultWidth:  stages.DefaultWidth,
		DefaultHeight: stages.DefaultHeight,
	})
	if err != nil {
		return Storie{}, err
	}

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

	videoFile, err := tpl.Render("video.mp4")
	if err != nil {
		return Storie{}, err
	}

	videoFile, err = stages.BuildVideo(ctx, stages.BuildVideoOptions{
		Stage:  stage,
		Target: videoFile,
	})

	return Storie{
		Stage: stage,
		Video: videoFile,
	}, nil
}
