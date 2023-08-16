package stages

import (
	"context"
	"os"

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

type Stage struct {
	Width      int
	Height     int
	Full       string
	Background string
	Foreground string
}

func BuildStage(ctx context.Context, opt BuildStageOptions) (Stage, error) {
	files := Stage{
		Width:  DefaultWidth,
		Height: DefaultHeight,
	}

	drawer, err := NewDraw(files.Width, files.Height)
	if err != nil {
		return files, err
	}

	buildStageImage := func(build drawPipe, name string) (string, error) {
		target, err := opt.Template.Render(name)
		if err != nil {
			return "", err
		}

		if err := build(opt.Source); err != nil {
			return "", ErrFailToBuildStage.Wrap(err)
		}

		targetFile, err := os.Create(target)
		if err != nil {
			return "", ErrFailToCreateFile.Wrap(err)
		}

		if err = drawer.Write(targetFile); err != nil {
			return "", ErrFailToWriteFile.Wrap(err)
		}

		return target, nil
	}

	if files.Full, err = buildStageImage(drawer.Draw, "full.png"); err != nil {
		return files, err
	}

	drawer.Reset()

	if files.Background, err = buildStageImage(drawer.DrawBase, "background.png"); err != nil {
		return files, err
	}

	drawer.Reset()

	if files.Foreground, err = buildStageImage(drawer.DrawOver, "foreground.png"); err != nil {
		return files, err
	}

	return files, err
}
