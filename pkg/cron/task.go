package cron

import (
	"context"

	"github.com/vinicius73/gamer-feed/pkg/model"
	"github.com/vinicius73/gamer-feed/pkg/tasks"
)

type TaskAction string

const (
	TaskSendLastEntries TaskAction = "send_last_entries"
)

var (
	_ tasks.Task[model.IEntry] = (*Task[model.IEntry, tasks.Task[model.IEntry]])(nil)
	_ tasks.Task[model.IEntry] = (*Task[model.IEntry, tasks.SendLastEntries[model.IEntry]])(nil)
	_ tasks.Task[model.IEntry] = (*Task[model.IEntry, tasks.Backup[model.IEntry]])(nil)
	_ tasks.Task[model.IEntry] = (*Task[model.IEntry, tasks.Cleanup[model.IEntry]])(nil)
)

type ScheduleTask[A model.IEntry] interface {
	tasks.Task[A]
	GetSchedules() []string
	Chats() []int64
}

type Task[A model.IEntry, T tasks.Task[A]] struct {
	Config    T        `fig:"config"    yaml:"config"`
	Schedules []string `fig:"schedules" yaml:"schedules"`
	ChatIDs   []int64  `fig:"chats"     yaml:"chats"`
}

func (t Task[A, T]) Name() string {
	return t.Config.Name()
}

func (t Task[A, T]) Run(ctx context.Context, opts tasks.TaskRunOptions[A]) error {
	return t.Config.Run(ctx, opts)
}

func (t Task[A, T]) GetSchedules() []string {
	return t.Schedules
}

func (t Task[A, T]) Chats() []int64 {
	return t.ChatIDs
}
