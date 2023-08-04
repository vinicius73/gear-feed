package configurations

import (
	"context"
)

type ctxKey struct{}

type AppConfig struct {
	Debug    bool   `fig:"-" yaml:"-"`
	Timezone string `fig:"timezone" yaml:"timezone"`
	Logger   Logger `fig:"logger" yaml:"logger"`
}

type Logger struct {
	Level  string `fig:"level" yaml:"level" default:"info"`
	Format string `fig:"format" yaml:"format" default:"text"`
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
