package tasks

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog"
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/sender"
)

var _ Task[model.IEntry] = (*Backup[model.IEntry])(nil)

type fileData struct {
	name string
	size int64
}

type Backup[T model.IEntry] struct {
	Base      string `fig:"base" yaml:"base"`
	Glob      string `fig:"glob" yaml:"glob"`
	AliasName string `fig:"name" yaml:"name"`
}

func (t Backup[T]) Name() string {
	return "backup"
}

func (t Backup[T]) Run(ctx context.Context, opts TaskRunOptions[T]) error {
	logger := zerolog.Ctx(ctx)

	tmpFile, err := os.CreateTemp(os.TempDir(), fmt.Sprintf("gfeed-backup--%s--*.tar", t.AliasName))
	if err != nil {
		return err
	}

	logger.Info().Str("file", tmpFile.Name()).Msg("tar file created")

	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	dataFiles, err := t.buildBackup(ctx, tmpFile)
	if err != nil {
		return err
	}

	err = opts.Sender.SendFile(ctx, sender.SendFileOptions{
		FilePath: tmpFile.Name(),
		Caption:  buildCaption(dataFiles),
	})

	if err != nil {
		return err
	}

	return nil
}

func (t Backup[T]) buildBackup(ctx context.Context, tmpFile *os.File) ([]fileData, error) {
	logger := zerolog.Ctx(ctx)

	glob := filepath.Join(t.Base, t.Glob)

	files, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}

	logger.Debug().Str("glob", glob).Int("files", len(files)).Msg("found files")

	tarWriter := tar.NewWriter(tmpFile)

	dataFiles := make([]fileData, len(files))

	for index, file := range files {
		dataFiles[index], err = addFileToTar(tarWriter, t.Base, file)

		if err != nil {
			return dataFiles, err
		}

		logger.Info().Str("file", file).Msg("adding file to tar")
	}

	if err := tarWriter.Close(); err != nil {
		return dataFiles, err
	}

	if err := tmpFile.Close(); err != nil {
		return dataFiles, err
	}

	logger.Info().Str("file", tmpFile.Name()).Msg("files stored in tar")

	return dataFiles, nil
}

func buildCaption(files []fileData) string {
	var capion strings.Builder

	capion.WriteString("<b>🕐 Backup </b>")
	capion.WriteString(time.Now().Format(time.RFC3339))
	capion.WriteRune('\n')
	capion.WriteString("🗃 Files: ")

	for _, file := range files {
		capion.WriteRune('\n')
		capion.WriteString("- <code>")
		capion.WriteString(file.name)
		capion.WriteString("</code> ")
		capion.WriteString(humanize.Bytes(uint64(file.size)))
	}

	return capion.String()
}

func addFileToTar(tarWriter *tar.Writer, base, file string) (fileData, error) {
	openedFile, err := os.Open(file)
	if err != nil {
		return fileData{}, err
	}

	defer openedFile.Close()

	stat, err := openedFile.Stat()
	if err != nil {
		return fileData{}, err
	}

	//nolint:exhaustruct
	hdr := &tar.Header{
		Name: strings.TrimPrefix(file, base),
		Mode: int64(stat.Mode()),
		Size: stat.Size(),
	}

	data := fileData{
		name: hdr.Name,
		size: hdr.Size,
	}

	if err := tarWriter.WriteHeader(hdr); err != nil {
		return data, err
	}

	if _, err := openedFile.Seek(0, 0); err != nil {
		return data, err
	}

	if _, err := io.Copy(tarWriter, openedFile); err != nil {
		return data, err
	}

	return data, nil
}
