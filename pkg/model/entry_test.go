package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vinicius73/gamer-feed/pkg/model"
)

func TestEntry_Hash(t *testing.T) {
	t.Parallel()

	type testFields struct {
		input model.Entry
		want  string
	}

	tests := []testFields{
		{
			input: model.Entry{
				Title:      "The Last of Us Part II",
				URL:        "https://www.gamereactor.eu/the-last-of-us-part-ii-review/",
				Image:      "https://images.gamereactor.eu/remote/articles/611893/The-Last-of-Us-Part-II-Review-0.jpg",
				Categories: []string{"review"},
				SourceName: "gamereactor",
			},
			want: "9311deeecac5fb039a8e3f6102659821f103ee48f486c55ce7c3868151ee25aa",
		},
		{
			input: model.Entry{
				Title:      "The Amazing Spider-Man",
				URL:        "https://www.gamereactor.eu/the-amazing-spiderman-review/",
				Image:      "https://images.gamereactor.eu/remote/articles/611893/The-Last-of-Us-Part-II-Review-0.jpg",
				Categories: []string{"review"},
				SourceName: "gamereactor",
			},
			want: "d3be2aa13ed6d3b659e5536016e1351e3e5ac2f85088e80e13524aa706c693e0",
		},
	}

	for _, test := range tests {
		entry := test.input
		got, err := entry.Hash()

		assert.Nil(t, err)
		assert.Equal(t, test.want, got)
	}
}
