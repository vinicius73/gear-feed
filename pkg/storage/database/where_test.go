package database_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vinicius73/gear-feed/pkg/model"
	"github.com/vinicius73/gear-feed/pkg/storage"
	"github.com/vinicius73/gear-feed/pkg/storage/database"
)

func TestWhere(t *testing.T) {
	t.Parallel()

	type testEntry struct {
		name  string
		where storage.WhereOptions
		item  database.DBEntry[model.Entry]
		has   bool
		want  bool
	}

	tests := []testEntry{
		{
			name:  "empty where and entry, also has false, should return false",
			where: storage.WhereOptions{},
			item:  database.DBEntry[model.Entry]{},
			has:   false,
			want:  false,
		},
		{
			name:  "empty where and entry, also has true, should return true",
			where: storage.WhereOptions{},
			item:  database.DBEntry[model.Entry]{},
			has:   true,
			want:  true,
		},
		{
			name:  "where Is new and empty item, should return false",
			where: storage.Where(storage.WhereIs(storage.StatusNew)),
			item:  database.DBEntry[model.Entry]{},
			has:   true,
			want:  false,
		},
		{
			name:  "where Is new and item is new, should return true",
			where: storage.Where(storage.WhereIs(storage.StatusNew)),
			item:  database.DBEntry[model.Entry]{Status: storage.StatusNew},
			has:   true,
			want:  true,
		},
		{
			name:  "where Is new and item is sent, should return false",
			where: storage.Where(storage.WhereIs(storage.StatusNew)),
			item:  database.DBEntry[model.Entry]{Status: storage.StatusSent},
			has:   true,
			want:  false,
		},
		{
			name:  "where Is sent and item is new, should return false",
			where: storage.Where(storage.WhereIs(storage.StatusSent)),
			item:  database.DBEntry[model.Entry]{Status: storage.StatusNew},
			has:   true,
			want:  false,
		},
		{
			name:  "where Is sent and item is sent, should return true",
			where: storage.Where(storage.WhereIs(storage.StatusSent)),
			item:  database.DBEntry[model.Entry]{Status: storage.StatusSent},
			has:   true,
			want:  true,
		},
		{
			name:  "where not new and item is new, should return false",
			where: storage.Where(storage.WhereNot(storage.StatusNew)),
			item:  database.DBEntry[model.Entry]{Status: storage.StatusNew},
			has:   true,
			want:  false,
		},
		{
			name:  "where not new and item is sent, should return true",
			where: storage.Where(storage.WhereNot(storage.StatusNew)),
			item:  database.DBEntry[model.Entry]{Status: storage.StatusSent},
			has:   true,
			want:  true,
		},
		{
			name:  "where not sent and item is new, should return true",
			where: storage.Where(storage.WhereNot(storage.StatusSent)),
			item:  database.DBEntry[model.Entry]{Status: storage.StatusNew},
			has:   true,
			want:  true,
		},
	}

	for _, tt := range tests {
		test := tt

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := database.Where(test.where, test.has, test.item)

			assert.Equal(t, test.want, got)
		})
	}
}
