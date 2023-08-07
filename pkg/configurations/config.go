package configurations

import (
	"context"

	"github.com/vinicius73/gamer-feed/pkg/cron"
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/storage/local"
	"github.com/vinicius73/gamer-feed/pkg/telegram"
)

type ctxKey struct{}

type AppConfig struct {
	Debug    bool                              `fig:"-"        yaml:"-"`
	Timezone string                            `fig:"timezone" yaml:"timezone"`
	Logger   Logger                            `fig:"logger"   yaml:"logger"`
	Telegram telegram.Config                   `fig:"telegram" yaml:"telegram"`
	Storage  local.Options                     `fig:"storage"  yaml:"storage"`
	Cron     cron.CronTasksConfig[model.Entry] `fig:"cron"     yaml:"cron"`
}

type Logger struct {
	Level  string `default:"info" fig:"level"  yaml:"level"`
	Format string `default:"text" fig:"format" yaml:"format"`
}

func (l Logger) Debug(level string) bool {
	debugLevels := [2]string{"debug", "trace"}

	for _, val := range debugLevels {
		if val == level {
			return true
		}
	}

	for _, val := range debugLevels {
		if val == l.Level {
			return true
		}
	}

	return false
}

func Ctx(ctx context.Context) *AppConfig {
	cf, _ := ctx.Value(ctxKey{}).(*AppConfig)

	return cf
}

func (c *AppConfig) WithContext(ctx context.Context) context.Context {
	if cf, ok := ctx.Value(ctxKey{}).(*AppConfig); ok {
		if cf == c {
			return ctx
		}
	}

	return context.WithValue(ctx, ctxKey{}, c)
}

func (c AppConfig) Tags() map[string]interface{} {
	return map[string]interface{}{}
}
