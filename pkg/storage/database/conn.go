package database

import (
	"context"
	"database/sql"
	"embed"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/vinicius73/gamer-feed/pkg/storage"
	"github.com/vinicius73/gamer-feed/pkg/support/apperrors"
	_ "modernc.org/sqlite" // Enable sqlite driver.
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

var (
	ErrFailToOpenDatabase = apperrors.System(nil, "fail to open database: %s", "DB:FAIL_TO_OPEN_DATABASE")
	ErrFailToRunMigration = apperrors.System(nil, "fail to run migration", "DB:FAIL_TO_RUN_MIGRATION")
)

type Options struct {
	storage.Options `fig:",squash" yaml:",inline"`
	Path            string `fig:"path"    yaml:"path"`
}

func Open(ctx context.Context, conf Options) (*sql.DB, error) {
	filename, err := checkDatabaseFile(conf.Path)
	if err != nil {
		return nil, ErrFailToOpenDatabase.Wrap(err).Msgf(filename)
	}

	conn, err := sql.Open("sqlite", filename)
	if err != nil {
		return nil, ErrFailToOpenDatabase.Wrap(err)
	}

	return conn, applyMigrations(ctx, conn)
}

func applyMigrations(ctx context.Context, conn *sql.DB) error {
	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: migrationsFS,
		Root:       "migrations",
	}

	count, err := migrate.ExecContext(ctx, conn, "sqlite3", migrations, migrate.Up)
	if err != nil {
		return ErrFailToRunMigration.Wrap(err)
	}

	logger := zerolog.Ctx(ctx).With().Str("context", "db:migrations").Logger()

	if count > 0 {
		logger.Info().Msgf("Applied migrations: %v", count)
	} else {
		logger.Info().Msg("No migrations to apply")
	}

	return nil
}

func checkDatabaseFile(filename string) (string, error) {
	if strings.HasSuffix(filename, ".sqlite") {
		return filename, nil
	}

	stat, err := os.Stat(filename)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	if stat.IsDir() {
		return filepath.Join(filename, "gfeed.sqlite"), nil
	}

	return filename, nil
}
