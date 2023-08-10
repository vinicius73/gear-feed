package sender_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vinicius73/gamer-feed/pkg/sender"
)

func TestCalculeSendInterval(t *testing.T) {
	t.Parallel()

	tests := []struct {
		count int
		want  time.Duration
	}{
		{
			count: 1,
			want:  time.Second,
		},
		{
			count: 2,
			want:  time.Second,
		},
		{
			count: 5,
			want:  time.Second,
		},
		{
			count: 9,
			want:  time.Second,
		},
		{
			count: 10,
			want:  time.Second * 2,
		},
		{
			count: 11,
			want:  time.Second * 2,
		},
		{
			count: 15,
			want:  time.Second * 2,
		},
		{
			count: 18,
			want:  time.Second * 2,
		},
		{
			count: 19,
			want:  time.Second * 2,
		},
		{
			count: 20,
			want:  time.Second * 3,
		},
		{
			count: 30,
			want:  time.Second * 3,
		},
		{
			count: 60,
			want:  time.Second * 3,
		},
	}

	for _, tt := range tests {
		actual := tt
		t.Run(fmt.Sprintf("CalculeSendInterval(%v) shold be %v", actual.count, actual.want), func(t *testing.T) {
			t.Parallel()
			got := sender.CalculeSendInterval(actual.count)
			assert.Equal(t, actual.want, got)
		})
	}
}
