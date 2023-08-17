package storage

import (
	"encoding/json"
	"time"

	"github.com/vinicius73/gamer-feed/pkg/model"
)

type Status byte

const (
	StatusNew Status = iota + 1
	StatusSent
)

type Options struct {
	TTL time.Duration `fig:"ttl" yaml:"ttl"`
}

type FindByHasStoryOptions struct {
	SourceNames []string
	Interval    time.Duration
	Limit       int
	Has         bool
}

type Entry[T model.IEntry] struct {
	Data   T
	Status Status
}

type Storage[T model.IEntry] interface {
	Has(hash string) (bool, error)
	Store(entry Entry[T]) error
	FindByHasStory(opt FindByHasStoryOptions) ([]T, error)
	Update(entry Entry[T]) error
	Cleanup() (int64, error)
	Where(opts WhereOptions, list []T) ([]T, error)
}

func (e Entry[T]) Hash() ([]byte, error) {
	hash, err := e.Data.Hash()
	if err != nil {
		return []byte{}, err
	}

	return []byte(hash), nil
}

func (e Entry[T]) Marshal() ([]byte, error) {
	return json.Marshal(e.Data)
}

func (s Status) Is(status Status) bool {
	return s == status
}

func (s Status) Byte() byte {
	return byte(s)
}
