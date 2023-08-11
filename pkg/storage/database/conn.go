package database

import (
	"context"
	"database/sql"
	"embed"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/vinicius73/gamer-feed/pkg/storage"
	"github.com/vinicius73/gamer-feed/pkg/support"
	"github.com/vinicius73/gamer-feed/pkg/support/apperrors"
	_ "modernc.org/sqlite" // Enable sqlite driver.
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

var (
	ErrFailToOpenDatabase = apperrors.System(nil, "fail to open database: %s", "DB:FAIL_TO_OPEN_DATABASE")
	ErrFailToRunMigration = apperrors.System(nil, "fail to run migration", "DB:FAIL_TO_RUN_MIGRATION")
	ErrDatabaseMustBeFile = apperrors.Business("database must be a file: %s", "DB:DATABASE_MUST_BE_FILE")
	ErrDatabaseMustExist  = apperrors.Business("database file must exist: %s", "DB:DATABASE_MUST_EXIST")
)

type Options struct {
	storage.Options `fig:",squash"    yaml:",inline"`
	Path            string `fig:"path"       yaml:"path"`
	MustExist       bool   `fig:"must_exist" yaml:"must_exist"`
}

func Open(ctx context.Context, conf Options) (*sql.DB, error) {
	logger := zerolog.Ctx(ctx).With().Str("context", "db:open").Str("filename", conf.Path).Logger()

	ctx = logger.WithContext(ctx)

	if conf.MustExist {
		if err := checkDatabaseFile(conf.Path); err != nil {
			return nil, err
		}
	} else if err := support.DirMustExist(filepath.Dir(conf.Path)); err != nil {
		return nil, err
	}

	conn, err := sql.Open("sqlite", conf.Path)
	if err != nil {
		return nil, ErrFailToOpenDatabase.Wrap(err).Msgf(conf.Path)
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

func checkDatabaseFile(filename string) error {
	stat, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return ErrDatabaseMustExist.Msgf(filename)
	}

	if stat.IsDir() {
		return ErrDatabaseMustBeFile.Msgf(filename)
	}

	return nil
}
