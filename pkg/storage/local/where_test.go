package local_test

import (
	"testing"

	"github.com/dgraph-io/badger/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vinicius73/gamer-feed/pkg/storage"
	"github.com/vinicius73/gamer-feed/pkg/storage/local"
)

var _ local.Record = testRecord{} // Ensure interface implementation

type testRecord struct {
	deletedOrExpired bool
	userMeta         byte
}

func (t testRecord) IsDeletedOrExpired() bool {
	return t.deletedOrExpired
}

func (t testRecord) UserMeta() byte {
	return t.userMeta
}

func TestWhere(t *testing.T) {
	t.Parallel()

	type testEntry struct {
		name  string
		where storage.WhereOptions
		item  testRecord
		rerr  error
		want  bool
	}

	tests := []testEntry{
		{
			name:  "empty where should return false",
			where: storage.WhereOptions{},
			item:  testRecord{},
			rerr:  nil,
			want:  false,
		},
		{
			name:  "where Is new and empty item should return false",
			where: storage.Where(storage.WhereIs(storage.StatusNew)),
			item:  testRecord{},
			rerr:  nil,
			want:  false,
		},
		{
			name:  "where Is new and item is new should return true",
			where: storage.Where(storage.WhereIs(storage.StatusNew)),
			item:  testRecord{userMeta: storage.StatusNew.Byte()},
			rerr:  nil,
			want:  true,
		},
		{
			name:  "where Is new and item is sent should return false",
			where: storage.Where(storage.WhereIs(storage.StatusNew)),
			item:  testRecord{userMeta: storage.StatusSent.Byte()},
			rerr:  nil,
			want:  false,
		},
		{
			name:  "where Not new and item is new should return false",
			where: storage.Where(storage.WhereNot(storage.StatusNew)),
			item:  testRecord{userMeta: storage.StatusNew.Byte()},
			rerr:  nil,
			want:  false,
		},
		{
			name:  "where Not new and item is sent should return true",
			where: storage.Where(storage.WhereNot(storage.StatusNew)),
			item:  testRecord{userMeta: storage.StatusSent.Byte()},
			rerr:  nil,
			want:  true,
		},
		{
			name:  "where Not new and item is new should return false",
			where: storage.Where(storage.WhereNot(storage.StatusNew)),
			item:  testRecord{userMeta: storage.StatusNew.Byte()},
			rerr:  nil,
			want:  false,
		},
		{
			name:  "where Not new and item expired return false",
			where: storage.Where(storage.WhereIs(storage.StatusNew)),
			item:  testRecord{deletedOrExpired: true, userMeta: storage.StatusNew.Byte()},
			rerr:  nil,
			want:  false,
		},
		{
			name:  "where Not new and item deleted return false",
			where: storage.Where(storage.WhereNot(storage.StatusNew)),
			item:  testRecord{deletedOrExpired: true, userMeta: storage.StatusSent.Byte()},
			rerr:  nil,
			want:  false,
		},
		{
			name:  "where allow missed shold return true",
			where: storage.Where(storage.WhereAllowMissed(true)),
			item:  testRecord{},
			rerr:  badger.ErrKeyNotFound,
			want:  true,
		},
		{
			name:  "where not allow missed shold return false",
			where: storage.Where(storage.WhereAllowMissed(false)),
			item:  testRecord{},
			rerr:  badger.ErrKeyNotFound,
			want:  false,
		},
	}

	for _, tt := range tests {
		current := tt
		t.Run(tt.name, func(t *testing.T) {
			got := local.ApplyWhere(current.where, current.item, current.rerr)

			assert.Equal(t, current.want, got)
		})
	}
}
