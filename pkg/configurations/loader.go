package configurations

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/creasty/defaults"
	"github.com/kkyr/fig"
	"github.com/mitchellh/go-homedir"
	"github.com/vinicius73/gamer-feed/pkg/support"
	"github.com/vinicius73/gamer-feed/pkg/support/apperrors"
	"gopkg.in/yaml.v3"
)

const (
	configBaseName = "gfeed"
	defaultTTL     = time.Hour * 30 * 24 // 30 days
)

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
		return fromFile(file)
	}

	home, err := homedir.Dir()
	if err != nil {
		return cfg, ErrFailToLoadConfig.Wrap(err)
	}

	err = fig.Load(&cfg,
		fig.File(configBaseName+".yml"),
		fig.UseEnv("GFEED"),
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

//nolint:cyclop,funlen
func applyDefaults(cfg AppConfig) (AppConfig, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return cfg, err
	}

	if cfg.Logger.Level == "" {
		cfg.Logger.Level = "info"
	}

	if cfg.Logger.Format == "" {
		cfg.Logger.Format = "text"
	}

	if cfg.Timezone == "" {
		loc, _ := time.LoadLocation("Local")

		cfg.Timezone = loc.String()
	}

	if cfg.Telegram.Token == "" {
		cfg.Telegram.Token = os.Getenv("TELEGRAM_TOKEN")
	}

	if cfg.Telegram.Token == "" {
		return cfg, ErrMissingTelegramToken
	}

	if cfg.Storage.Path == "" {
		cfg.Storage.Path = support.GetEnvString("GFEED_DATABASE", "."+configBaseName+".sqlite")
	}

	if !path.IsAbs(cfg.Storage.Path) {
		cfg.Storage.Path = path.Join(pwd, cfg.Storage.Path)
	}

	cfg.Storage.Path, err = checkDatabaseFile(cfg.Storage.Path)

	if err != nil {
		return cfg, err
	}

	if cfg.Storage.TTL == 0 {
		cfg.Storage.TTL = defaultTTL
	}

	cfg.Cron.Timezone, _ = time.LoadLocation(cfg.Timezone)

	if cfg.Cron.Backup.Config.Base != "" && !filepath.IsAbs(cfg.Cron.Backup.Config.Base) {
		cfg.Cron.Backup.Config.Base = path.Join(pwd, cfg.Cron.Backup.Config.Base)
	}

	if cfg.Cron.Backup.Config.AliasName == "" {
		hostname, err := os.Hostname()
		if err != nil {
			return cfg, err
		}

		cfg.Cron.Backup.Config.AliasName = hostname
	}

	return cfg, nil
}

func ensureConfig() (AppConfig, error) {
	var err error

	var cfg AppConfig

	if err = defaults.Set(&cfg); err != nil {
		return cfg, ErrFailEnsureConfig.Wrap(err)
	}

	cfg, err = applyDefaults(cfg)

	if err != nil {
		return cfg, ErrFailEnsureConfig.Wrap(err)
	}

	buf, err := yaml.Marshal(cfg)
	if err != nil {
		return cfg, ErrFailEnsureConfig.Wrap(err)
	}

	pwd, _ := os.Getwd()

	configFile := path.Join(pwd, configBaseName+".yml")

	err = os.WriteFile(configFile, buf, os.ModePerm)

	if err != nil {
		return cfg, ErrFailEnsureConfig.Wrap(err)
	}

	return cfg, ConfigFileWasCreated.Msgf(configFile)
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
		return filepath.Join(filename, configBaseName+".sqlite"), nil
	}

	return filename, nil
}

func fromFile(file string) (AppConfig, error) {
	var cfg AppConfig

	content, err := loadFileContentAndApplyEnv(file)

	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(content, &cfg)

	if err != nil {
		return cfg, err
	}

	return applyDefaults(cfg)
}

func loadFileContentAndApplyEnv(file string) ([]byte, error) {
	content, err := os.ReadFile(file)

	if err != nil {
		return nil, err
	}

	return []byte(os.ExpandEnv(string(content))), nil
}
