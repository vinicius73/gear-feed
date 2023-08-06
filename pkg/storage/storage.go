package storage

import (
	"encoding/json"
	"time"
)

type Status byte

const (
	StatusNew Status = iota + 1
	StatusSent
)

type Options struct {
	TTL time.Duration
}

type Hashable interface {
	Hash() (string, error)
}

type Entry[T Hashable] struct {
	Data   T
	Status Status
}

type Storage[T Hashable] interface {
	Has(hash string) (bool, error)
	Store(entry Entry[T]) error
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
