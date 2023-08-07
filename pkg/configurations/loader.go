package configurations

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/creasty/defaults"
	"github.com/kkyr/fig"
	"github.com/mitchellh/go-homedir"
	"github.com/vinicius73/gamer-feed/pkg/support"
	"github.com/vinicius73/gamer-feed/pkg/support/apperrors"
	"gopkg.in/yaml.v3"
)

const configBaseName = "gfeed"

var (
	ConfigFileWasCreated    = apperrors.Business("a new config file was created (%s)", "CONF:001")
	ErrFailEnsureConfig     = apperrors.System(nil, "fail to ensure config", "CONF:002")
	ErrMissingTelegramToken = apperrors.System(nil, "missing telegram token", "CONF:003")
	ErrFailToLoadConfig     = apperrors.System(nil, "fail to load config", "CONF:004")
)

func Load(file string) (AppConfig, error) {
	var err error

	var cfg AppConfig

	if file != "" {
		err = fig.Load(&cfg,
			fig.File(filepath.Base(file)),
			fig.Dirs(filepath.Dir(file)),
		)

		if err != nil {
			return cfg, ErrFailToLoadConfig.WithErr(err)
		}

		return applyDefaults(cfg)
	}

	home, err := homedir.Dir()
	if err != nil {
		return cfg, ErrFailToLoadConfig.WithErr(err)
	}

	err = fig.Load(&cfg,
		fig.File(configBaseName+".yml"),
		fig.Dirs(
			".",
			path.Join(home, "."+configBaseName),
			path.Join(home, ".config"),
			path.Join(home, ".config/"+configBaseName),
			home,
			"/etc/"+configBaseName,
			"/"+configBaseName+".d",
			support.GetBinDirPath(),
		),
	)

	if errors.Is(err, fig.ErrFileNotFound) {
		return ensureConfig()
	}

	if err != nil {
		return cfg, err
	}

	return applyDefaults(cfg)
}

func applyDefaults(cfg AppConfig) (AppConfig, error) {
	if cfg.Logger.Level == "" {
		cfg.Logger.Level = "info"
	}

	if cfg.Logger.Format == "" {
		cfg.Logger.Format = "text"
	}

	if cfg.Timezone == "" {
		cfg.Timezone = time.Local.String()
	}

	if cfg.Telegram.Token == "" {
		cfg.Telegram.Token = os.Getenv("TELEGRAM_TOKEN")
	}

	if cfg.Telegram.Token == "" {
		return cfg, ErrMissingTelegramToken
	}

	if cfg.Storage.Path == "" {
		cfg.Storage.Path = "./." + configBaseName + ".db"
	}

	if path.IsAbs(cfg.Storage.Path) {
		pwd, _ := os.Getwd()
		cfg.Storage.Path = path.Join(pwd, cfg.Storage.Path)
	}

	if cfg.Storage.TTL == 0 {
		cfg.Storage.TTL = 30 * 24 * time.Hour
	}

	cfg.Cron.Timezone, _ = time.LoadLocation(cfg.Timezone)

	return cfg, nil
}

func ensureConfig() (AppConfig, error) {
	var err error

	var cfg AppConfig

	if err = defaults.Set(&cfg); err != nil {
		return cfg, ErrFailEnsureConfig.WithErr(err)
	}

	cfg, err = applyDefaults(cfg)

	if err != nil {
		return cfg, ErrFailEnsureConfig.WithErr(err)
	}

	buf, err := yaml.Marshal(cfg)
	if err != nil {
		return cfg, ErrFailEnsureConfig.WithErr(err)
	}

	pwd, _ := os.Getwd()

	configFile := path.Join(pwd, configBaseName+".yml")

	err = os.WriteFile(configFile, buf, os.ModePerm)

	if err != nil {
		return cfg, ErrFailEnsureConfig.WithErr(err)
	}

	return cfg, ConfigFileWasCreated.Msgf(configFile)
}
