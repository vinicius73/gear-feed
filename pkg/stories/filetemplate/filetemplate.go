package filetemplate

import (
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gosimple/slug"
	"github.com/vinicius73/gear-feed/pkg/stories/fetcher"
	"github.com/vinicius73/gear-feed/pkg/support"
	"github.com/vinicius73/gear-feed/pkg/support/apperrors"
)

const (
	maxTitleName = 10
	hashSize     = 8
)

var (
	ErrEmptyDir            = apperrors.Business("base dir cannot be empty", "STAGES:EMPTY_TARGET")
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
	Hash    string
}

type Options struct {
	Source   fetcher.Result
	Template string
	BaseDir  string
}

func New(opt Options) (Template, error) {
	if opt.BaseDir == "" {
		return Template{}, ErrEmptyDir
	}

	if err := support.DirMustExist(opt.BaseDir); err != nil {
		return Template{}, ErrFailtoCreateDir.Wrap(err)
	}

	tpl, err := template.New("").Parse(opt.Template)
	if err != nil {
		return Template{}, ErrFailToParseTemplate.Wrap(err)
	}

	hash := opt.Source.Hash

	if len(hash) > hashSize {
		hash = hash[:hashSize]
	}

	title := opt.Source.Title

	if len(title) > maxTitleName {
		title = title[:maxTitleName]
	}

	return Template{
		tpl:     tpl,
		Hash:    hash,
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
		"hash":     t.Hash,
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
