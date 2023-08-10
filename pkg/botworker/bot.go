package botworker

import (
	"context"
	"time"

	"github.com/vinicius73/gamer-feed/pkg/cron"
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/sender"
	"github.com/vinicius73/gamer-feed/pkg/storage"
)

const stopTimeout = time.Second * 30

type Config[T model.IEntry] struct {
	Cron cron.TasksConfig[T]
}

type BotOptions[T model.IEntry] struct {
	Config  Config[T]
	Sender  sender.Serder[T]
	Storage storage.Storage[T]
}

type Bot[T model.IEntry] struct {
	config  Config[T]
	sender  sender.Serder[T]
	storage storage.Storage[T]
}

func New[T model.IEntry](opts BotOptions[T]) Bot[T] {
	return Bot[T]{
		config:  opts.Config,
		sender:  opts.Sender,
		storage: opts.Storage,
	}
}

func (b Bot[T]) Run(ctx context.Context) error {
	runner, err := b.schedule(ctx)
	if err != nil {
		return err
	}

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), stopTimeout)

	defer cancel()

	//nolint:contextcheck
	return runner.Stop(ctx)
}

func (b Bot[T]) schedule(ctx context.Context) (cron.Runner[T], error) {
	runner := cron.New[T](cron.RunnerOptions[T]{
		Storage: b.storage,
		Sender:  b.sender,
		Config:  b.config.Cron,
	})

	if err := runner.Start(ctx); err != nil {
		return runner, err
	}

	return runner, nil
}
