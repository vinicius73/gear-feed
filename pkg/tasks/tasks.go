package tasks

import (
	"context"

	"github.com/vinicius73/gear-feed/pkg/model"
	"github.com/vinicius73/gear-feed/pkg/sender"
	"github.com/vinicius73/gear-feed/pkg/storage"
)

type TaskRunOptions[T model.IEntry] struct {
	Storage storage.Storage[T]
	Sender  sender.Serder[T]
}

type Task[T model.IEntry] interface {
	Name() string
	Run(context.Context, TaskRunOptions[T]) error
}
