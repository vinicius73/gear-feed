package stages

import (
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gosimple/slug"
	"github.com/vinicius73/gamer-feed/pkg/support"
	"github.com/vinicius73/gamer-feed/pkg/support/apperrors"
)

var (
	ErrEmptyTarget = apperrors.Business("target cannot be empty", "COVERS:EMPTY_TARGET")
)

type TemplateFilename struct {
	Template *template.Template
	BaseDir  string
	Date     string
	Site     string
	Title    string
}

type BuildStageOptions struct {
	Source           FetchResult
	TargetDir        string
	TemplateFilename string
}

func (t TemplateFilename) Render(filename string) (string, error) {
	var builder strings.Builder

	err := t.Template.Execute(&builder, map[string]any{
		"baseDir":  t.BaseDir,
		"date":     t.Date,
		"site":     t.Site,
		"title":    t.Title,
		"filename": filename,
	})
	if err != nil {
		return "", err
	}

	result := builder.String()

	if strings.HasPrefix(result, t.BaseDir) {
		return result, nil
	}

	return filepath.Join(t.BaseDir, result), nil
}

func (bo BuildStageOptions) Template(source FetchResult) (TemplateFilename, error) {
	if bo.TargetDir == "" {
		return TemplateFilename{}, ErrEmptyTarget
	}

	if err := support.DirMustExist(filepath.Dir(bo.TargetDir)); err != nil {
		return TemplateFilename{}, err
	}

	tpl, err := template.New("").Parse(bo.TemplateFilename)
	if err != nil {
		return TemplateFilename{}, err
	}

	fileTemplate := TemplateFilename{
		Template: tpl,
		BaseDir:  bo.TargetDir,
		Date:     strconv.FormatInt(time.Now().Unix(), 10),
		Title:    slug.Make(source.Title),
		Site:     slug.Make(source.SiteName),
	}

	return fileTemplate, nil
}
