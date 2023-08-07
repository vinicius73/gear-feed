package cron

import (
	"context"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gosimple/slug"
	"github.com/rs/zerolog"
	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/sender"
	"github.com/vinicius73/gamer-feed/pkg/storage"
	"github.com/vinicius73/gamer-feed/pkg/support"
	"github.com/vinicius73/gamer-feed/pkg/tasks"
)

type TasksConfig[T model.IEntry] struct {
	Timezone        *time.Location                    `fig:"-"                 yaml:"-"`
	SendLastEntries Task[T, tasks.SendLastEntries[T]] `fig:"send_last_entries" yaml:"send_last_entries"`
}

type Runner[T model.IEntry] struct {
	storage   storage.Storage[T]
	sender    sender.Serder[T]
	config    TasksConfig[T]
	scheduler *gocron.Scheduler
}

type RunnerOptions[T model.IEntry] struct {
	Config  TasksConfig[T]
	Storage storage.Storage[T]
	Sender  sender.Serder[T]
}

func New[T model.IEntry](opts RunnerOptions[T]) Runner[T] {
	scheduler := gocron.NewScheduler(opts.Config.Timezone)
	scheduler.SingletonModeAll()
	scheduler.SetMaxConcurrentJobs(1, gocron.RescheduleMode)

	return Runner[T]{
		config:    opts.Config,
		storage:   opts.Storage,
		sender:    opts.Sender,
		scheduler: scheduler,
	}
}

func (r Runner[T]) Start(ctx context.Context) error {
	r.scheduler.Clear()

	logger := zerolog.Ctx(ctx)

	logger.Info().Msg("Starting cron tasks")

	tasks := []Task[T, tasks.Task[T]]{
		{
			Config:    r.config.SendLastEntries.Config,
			ChatIDs:   r.config.SendLastEntries.ChatIDs,
			Schedules: r.config.SendLastEntries.Schedules,
		},
	}

	for _, task := range tasks {
		err := r.register(ctx, task)
		if err != nil {
			return err
		}
	}

	r.scheduler.StartAsync()

	logger.Info().Msg("Cron tasks started")

	return nil
}

func (r Runner[T]) Stop(_ context.Context) error {
	r.scheduler.Clear()
	r.scheduler.Stop()

	return nil
}

func (r Runner[T]) register(ctx context.Context, task Task[T, tasks.Task[T]]) error {
	logger := zerolog.Ctx(ctx).With().Str("task", task.Config.Name()).Logger()

	if len(task.Schedules) == 0 {
		logger.Warn().Msg("Task has no schedules")

		return nil
	}

	for _, schedule := range task.Schedules {
		job, err := r.scheduler.Cron(schedule).Do(r.exec, ctx, task)
		if err != nil {
			return err
		}

		hash, err := support.HashSHA256(schedule)
		if err != nil {
			return err
		}

		name := slug.Make(task.Config.Name() + "_" + hash)

		job.Name(name)
		job.Tag(task.Config.Name())

		logger.Info().
			Str("task", task.Config.Name()).
			Str("schedule", schedule).
			Msg("Task registered")
	}

	return nil
}

func (r Runner[T]) exec(ctx context.Context, task Task[T, tasks.Task[T]]) {
	logger := zerolog.Ctx(ctx).With().Str("task", task.Config.Name()).Logger()
	ctx = logger.WithContext(ctx)

	logger.Info().Msg("Running task")

	opts := tasks.TaskRunOptions[T]{
		Storage: r.storage,
		Sender:  r.sender.WithChats(task.ChatIDs),
	}

	err := task.Run(ctx, opts)
	if err != nil {
		logger.
			Error().
			Err(err).
			Msg("Error running task")
	}

	logger.Info().Msg("Task finished")
}
