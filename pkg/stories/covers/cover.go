package covers

import (
	"os"
	"path/filepath"
	"strconv"
	"text/template"
	"time"

	"github.com/gosimple/slug"
	"github.com/vinicius73/gamer-feed/pkg/support"
	"github.com/vinicius73/gamer-feed/pkg/support/apperrors"
)

const (
	defaultWidth  = 1080
	defaultHeight = 1920
)

var (
	ErrEmptyTarget   = apperrors.Business("target cannot be empty", "COVERS:EMPTY_TARGET")
	ErrInvalidTarget = apperrors.Business("target must be a .png file", "COVERS:INVALID_TARGET")
)

type BuildOptions struct {
	TargetDir        string
	TemplateFilename string
}

type CoverBuilder struct {
	Source FetchResult
	Width  int
	Height int
}

func NewBuilder(source FetchResult) CoverBuilder {
	return CoverBuilder{
		Source: source,
		Width:  defaultWidth,
		Height: defaultHeight,
	}
}

func (b CoverBuilder) filePath(target, templateSource string) (TemplateFilename, error) {
	if target == "" {
		return TemplateFilename{}, ErrEmptyTarget
	}

	if err := support.DirMustExist(filepath.Dir(target)); err != nil {
		return TemplateFilename{}, err
	}

	tpl, err := template.New("filename").Parse(templateSource)
	if err != nil {
		return TemplateFilename{}, err
	}

	fileTemplate := TemplateFilename{
		Template: tpl,
		BaseDir:  target,
		Date:     strconv.FormatInt(time.Now().Unix(), 10),
		Title:    slug.Make(b.Source.Title),
		Site:     slug.Make(b.Source.SiteName),
	}

	return fileTemplate, nil
}

func (b CoverBuilder) Build(opt BuildOptions) ([]string, error) {
	drawer, err := NewDraw(b.Width, b.Height)
	files := make([]string, 0, 3)

	if err != nil {
		return nil, err
	}

	target, err := b.filePath(opt.TargetDir, opt.TemplateFilename)
	if err != nil {
		return nil, err
	}

	full, err := b.buildFull(drawer, target)
	if err != nil {
		return nil, err
	}

	files = append(files, full)

	base, err := b.buildBase(drawer, target)
	if err != nil {
		return nil, err
	}

	files = append(files, base)

	over, err := b.buildOver(drawer, target)
	if err != nil {
		return nil, err
	}

	files = append(files, over)

	return files, nil
}

func (b CoverBuilder) buildFull(drawer *Draw, tpl TemplateFilename) (string, error) {
	drawer.Reset()
	return b.build(drawer, drawer.Draw, "full.png", tpl)
}

func (b CoverBuilder) buildBase(drawer *Draw, tpl TemplateFilename) (string, error) {
	drawer.Reset()
	return b.build(drawer, drawer.DrawBase, "base.png", tpl)
}

func (b CoverBuilder) buildOver(drawer *Draw, tpl TemplateFilename) (string, error) {
	drawer.Reset()
	return b.build(drawer, drawer.DrawOver, "over.png", tpl)
}

func (b CoverBuilder) build(drawer *Draw, build drawPipe, name string, tpl TemplateFilename) (string, error) {
	target, err := tpl.Render(name)
	if err != nil {
		return "", err
	}

	if err := build(b.Source); err != nil {
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
