package filetemplate

import (
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gosimple/slug"
	"github.com/vinicius73/gamer-feed/pkg/stories/fetcher"
	"github.com/vinicius73/gamer-feed/pkg/support"
	"github.com/vinicius73/gamer-feed/pkg/support/apperrors"
)

const maxTitleName = 10

var (
	ErrEmptyTarget         = apperrors.Business("target cannot be empty", "STAGES:EMPTY_TARGET")
	ErrFailtoCreateDir     = apperrors.System(nil, "fail to create dir", "STAGES:FAIL_TO_CREATE_DIR")
	ErrFailToParseTemplate = apperrors.System(nil, "fail to parse template", "STAGES:FAIL_TO_PARSE_TEMPLATE")
	ErrFailtoBuildFilename = apperrors.System(nil, "fail to build filename", "STAGES:FAIL_TO_BUILD_FILENAME")
)

type Template struct {
	tpl     *template.Template
	BaseDir string
	Date    string
	Site    string
	Title   string
}

type Options struct {
	Source   fetcher.Result
	Template string
	BaseDir  string
}

func New(opt Options) (Template, error) {
	tpl, err := template.New("").Parse(opt.Template)
	if err != nil {
		//nolint:exhaustruct
		return Template{}, ErrFailToParseTemplate.Wrap(err)
	}

	if err = support.DirMustExist(opt.BaseDir); err != nil {
		//nolint:exhaustruct
		return Template{}, ErrFailtoCreateDir.Wrap(err)
	}

	title := opt.Source.Title

	if len(title) > maxTitleName {
		title = title[:maxTitleName]
	}

	return Template{
		tpl:     tpl,
		BaseDir: opt.BaseDir,
		Site:    slug.Make(opt.Source.SiteName),
		Title:   slug.Make(title),
		Date:    strconv.Itoa(int(time.Now().Unix())),
	}, nil
}

func (t Template) Render(filename string) (string, error) {
	var builder strings.Builder

	err := t.tpl.Execute(&builder, map[string]any{
		"baseDir":  t.BaseDir,
		"date":     t.Date,
		"site":     t.Site,
		"title":    t.Title,
		"filename": filename,
	})
	if err != nil {
		return "", ErrFailtoBuildFilename.Wrap(err)
	}

	result := builder.String()

	if strings.HasPrefix(result, t.BaseDir) {
		return result, nil
	}

	return filepath.Join(t.BaseDir, result), nil
}
