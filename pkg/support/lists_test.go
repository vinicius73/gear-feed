package support_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vinicius73/gamer-feed/pkg/support"
)

func TestFindIndex(t *testing.T) {
	type testDefiniton[T any] struct {
		list  []T
		test  func(T) bool
		index int
		found bool
	}

	ints := []testDefiniton[int]{
		{[]int{1, 2, 3}, func(i int) bool { return i == 2 }, 1, true},
		{[]int{1, 2, 3}, func(i int) bool { return i == 4 }, -1, false},
		{[]int{1, 2, 3}, func(i int) bool { return i == 1 }, 0, true},
		{[]int{1, 2, 3}, func(i int) bool { return i == 3 }, 2, true},
		{[]int{}, func(i int) bool { return i == 3 }, -1, false},
		{[]int{1, 2, 3}, func(i int) bool { return i == 0 }, -1, false},
	}

	for _, test := range ints {
		index, found := support.FindIndex(test.list, test.test)

		assert.Equal(t, test.index, index)
		assert.Equal(t, test.found, found)
	}

	strings := []testDefiniton[string]{
		{[]string{"a", "b", "c"}, func(i string) bool { return i == "b" }, 1, true},
		{[]string{"a", "b", "c"}, func(i string) bool { return i == "d" }, -1, false},
		{[]string{"a", "b", "c"}, func(i string) bool { return i == "a" }, 0, true},
		{[]string{"a", "b", "c"}, func(i string) bool { return i == "c" }, 2, true},
		{[]string{}, func(i string) bool { return i == "c" }, -1, false},
		{[]string{"a", "b", "c"}, func(i string) bool { return i == "e" }, -1, false},
	}

	for _, test := range strings {
		index, found := support.FindIndex(test.list, test.test)

		assert.Equal(t, test.index, index)
		assert.Equal(t, test.found, found)
	}
}

func TestContains(t *testing.T) {
	type testDefiniton[T comparable] struct {
		list []T
		val  T
		ok   bool
	}

	ints := []testDefiniton[int]{
		{[]int{1, 2, 3}, 2, true},
		{[]int{1, 2, 3}, 4, false},
		{[]int{1, 2, 3}, 1, true},
		{[]int{1, 2, 3}, 3, true},
		{[]int{}, 3, false},
		{[]int{1, 2, 3}, 0, false},
	}

	for _, test := range ints {
		assert.Equal(t, test.ok, support.Contains(test.list, test.val))
	}

	strings := []testDefiniton[string]{
		{[]string{"a", "b", "c"}, "b", true},
		{[]string{"a", "b", "c"}, "d", false},
		{[]string{"a", "b", "c"}, "a", true},
		{[]string{"a", "b", "c"}, "c", true},
		{[]string{}, "c", false},
		{[]string{"a", "b", "c"}, "e", false},
	}

	for _, test := range strings {
		assert.Equal(t, test.ok, support.Contains(test.list, test.val))
	}
}

func TestContainsSome(t *testing.T) {
	type testDefiniton[T comparable] struct {
		list []T
		vals []T
		ok   bool
	}

	ints := []testDefiniton[int]{
		{[]int{1, 2, 3}, []int{2, 3}, true},
		{[]int{1, 2, 3}, []int{4}, false},
		{[]int{1, 2, 3}, []int{1}, true},
		{[]int{1, 2, 3}, []int{3}, true},
		{[]int{}, []int{3}, false},
		{[]int{1, 2, 3}, []int{0}, false},
	}

	for _, test := range ints {
		assert.Equal(t, test.ok, support.ContainsSome(test.list, test.vals))
	}

	strings := []testDefiniton[string]{
		{[]string{"a", "b", "c"}, []string{"b", "c"}, true},
		{[]string{"a", "b", "c"}, []string{"d"}, false},
		{[]string{"a", "b", "c"}, []string{"a"}, true},
		{[]string{"a", "b", "c"}, []string{"c"}, true},
		{[]string{}, []string{"c"}, false},
		{[]string{"a", "b", "c"}, []string{"e"}, false},
	}

	for _, test := range strings {
		assert.Equal(t, test.ok, support.ContainsSome(test.list, test.vals))
	}
}

func TestToLower(t *testing.T) {
	type testDefiniton struct {
		list   []string
		result []string
	}

	tests := []testDefiniton{
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{[]string{"A", "B", "C"}, []string{"a", "b", "c"}},
		{[]string{"A", "b", "C"}, []string{"a", "b", "c"}},
		{[]string{"a", "B", "c"}, []string{"a", "b", "c"}},
		{[]string{"a", "b", "C"}, []string{"a", "b", "c"}},
		{[]string{"A", "B", "c"}, []string{"a", "b", "c"}},
		{[]string{"a", "B", "C"}, []string{"a", "b", "c"}},
		{[]string{"A", "b", "c"}, []string{"a", "b", "c"}},
	}

	for _, test := range tests {
		assert.Equal(t, test.result, support.ToLower(test.list))
	}
}
