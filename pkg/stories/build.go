package stories

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/vinicius73/gamer-feed/pkg/stories/fetcher"
	"github.com/vinicius73/gamer-feed/pkg/stories/filetemplate"
	"github.com/vinicius73/gamer-feed/pkg/stories/stages"
	"github.com/vinicius73/gamer-feed/pkg/support"
)

const workerSize = 2

type Footer stages.Footer

type BuildStorieOptions struct {
	Footer           Footer
	TemplateFilename string
	SourceURL        string
	TargetDir        string
}

type BuildCollectionOptions struct {
	Footer           Footer
	Sources          []string
	TemplateFilename string
	TargetDir        string
}

type workerResult struct {
	Story Story
	Error error
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

func BuildStory(ctx context.Context, opt BuildStorieOptions) (Story, error) {
	logger := zerolog.Ctx(ctx).With().Str("component", "stories").Logger()
	ctx = logger.WithContext(ctx)

	entry, err := fetcher.Fetch(ctx, fetcher.Options{
		SourceURL:     opt.SourceURL,
		DefaultWidth:  stages.DefaultWidth,
		DefaultHeight: stages.DefaultHeight,
	})
	if err != nil {
		return Story{}, err
	}

	logger = logger.With().Str("title", entry.Title).Logger()
	ctx = logger.WithContext(ctx)

	logger.Info().Msg("entry data loades")

	tpl, err := opt.Template(entry)
	if err != nil {
		return Story{}, err
	}

	stage, err := stages.BuildStage(ctx, stages.BuildStageOptions{
		Source:   entry,
		Template: tpl,
		Footer:   stages.Footer(opt.Footer),
	})
	if err != nil {
		return Story{}, err
	}

	logger.Info().Msg("stage builded")

	videoFile, err := tpl.Render("video.mp4")
	if err != nil {
		return Story{}, err
	}

	videoFile, err = stages.BuildVideo(ctx, stages.BuildVideoOptions{
		Stage:  stage,
		Target: videoFile,
	})
	if err != nil {
		return Story{}, err
	}

	logger.Info().Msg("video builded")

	return Story{
		Stage: stage,
		Video: videoFile,
		Hash:  entry.Hash,
	}, nil
}

func BuildCollection(ctx context.Context, opt BuildCollectionOptions) (Collection, error) {
	logger := zerolog.Ctx(ctx).With().Str("component", "stories").Logger()
	ctx = logger.WithContext(ctx)

	var collection Collection

	//nolint:gomnd
	input := make(chan BuildStorieOptions, 2)
	outs := make([]<-chan workerResult, workerSize)

	for i := range workerSize {
		outs[i] = buildWorker(ctx, input)
	}

	for _, source := range opt.Sources {
		input <- BuildStorieOptions{
			SourceURL:        source,
			TargetDir:        opt.TargetDir,
			TemplateFilename: opt.TemplateFilename,
			Footer:           opt.Footer,
		}
	}

	close(input)

	for res := range support.MergeChanners(outs...) {
		if res.Error != nil {
			logger.Error().Err(res.Error).Msg("Error on build worker")

			continue
		}

		collection = append(collection, res.Story)
	}

	return collection, nil
}

func buildWorker(ctx context.Context, input <-chan BuildStorieOptions) <-chan workerResult {
	//nolint:gomnd
	out := make(chan workerResult, 2)

	go func() {
		defer close(out)

		for {
			select {
			case opt, ok := <-input:
				if !ok {
					return
				}

				story, err := BuildStory(ctx, opt)
				out <- workerResult{Story: story, Error: err}
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}
