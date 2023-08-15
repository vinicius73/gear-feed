package covers

import (
	"path/filepath"
	"strings"
	"text/template"
)

type TemplateFilename struct {
	Template *template.Template
	BaseDir  string
	Date     string
	Site     string
	Title    string
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
