package stages

import "os"

const (
	defaultWidth  = 1080
	defaultHeight = 1920
)

type Stage struct {
	Width      int
	Height     int
	Full       string
	Background string
	Foreground string
}

func BuildStage(opt BuildStageOptions) (Stage, error) {
	files := Stage{

		Width:  defaultWidth,
		Height: defaultHeight,
	}

	drawer, err := NewDraw(files.Width, files.Height)

	if err != nil {
		return files, err
	}

	tpl, err := opt.Template(opt.Source)
	if err != nil {
		return files, err
	}

	buildStageImage := func(build drawPipe, name string) (string, error) {
		target, err := tpl.Render(name)
		if err != nil {
			return "", err
		}

		if err := build(opt.Source); err != nil {
			return "", err
		}

		targetFile, err := os.Create(target)
		if err != nil {
			return "", err
		}

		if err = drawer.Write(targetFile); err != nil {
			return "", err
		}

		return target, nil
	}

	if files.Full, err = buildStageImage(drawer.Draw, "full.png"); err != nil {
		return files, err
	}

	if files.Background, err = buildStageImage(drawer.DrawBase, "background.png"); err != nil {
		return files, err
	}

	if files.Foreground, err = buildStageImage(drawer.DrawOver, "foreground.png"); err != nil {
		return files, err
	}

	_, err = BuildVideo(BuildVideoOptions{
		Stage: files,
	})

	return files, err
}
