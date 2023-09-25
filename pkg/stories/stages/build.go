package stages

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/vinicius73/gamer-feed/pkg/stories/drawer"
	"github.com/vinicius73/gamer-feed/pkg/support/apperrors"
)

var (
	ErrFailToCreateFile = apperrors.System(nil, "fail to create file", "STAGES:FAIL_TO_CREATE_FILE")
	ErrFailToWriteFile  = apperrors.System(nil, "fail to write file", "STAGES:FAIL_TO_WRITE_FILE")
	ErrFailToBuildStage = apperrors.System(nil, "fail to build stage", "STAGES:FAIL_TO_BUILD_STAGE")
)

const (
	DefaultWidth  = 1080
	DefaultHeight = 1920
)

func BuildStage(ctx context.Context, opt BuildStageOptions) (Stage, error) {
	files := Stage{
		Width:      DefaultWidth,
		Height:     DefaultHeight,
		Full:       "",
		Background: "",
		Foreground: "",
	}

	draw, err := drawer.NewDraw(drawer.DrawOptions{
		Width:  files.Width,
		Height: files.Height,
		Footer: drawer.Footer(opt.Footer),
	})
	if err != nil {
		return files, err
	}

	buildStageImage := func(build drawer.DrawPipe, name string) (string, error) {
		logger := zerolog.Ctx(ctx).With().Str("stage", name).Logger()

		logger.Debug().Msg("building stage")

		target, err := opt.Template.Render(name)
		if err != nil {
			return "", err
		}

		if err := build(logger.WithContext(ctx), opt.Source); err != nil {
			return "", ErrFailToBuildStage.Wrap(err)
		}

		targetFile, err := os.Create(target)
		if err != nil {
			return "", ErrFailToCreateFile.Wrap(err)
		}

		if err = draw.Write(targetFile); err != nil {
			return "", ErrFailToWriteFile.Wrap(err)
		}

		return target, nil
	}

	if files.Full, err = buildStageImage(draw.Draw, "full.png"); err != nil {
		return files, err
	}

	draw.Reset()

	if files.Background, err = buildStageImage(draw.DrawBase, "background.png"); err != nil {
		return files, err
	}

	draw.Reset()

	if files.Foreground, err = buildStageImage(draw.DrawOver, "foreground.png"); err != nil {
		return files, err
	}

	return files, err
}
