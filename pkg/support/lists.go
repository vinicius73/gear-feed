package support

import "strings"

func FindIndex[T any](list []T, test func(T) bool) (int, bool) {
	for i, item := range list {
		if test(item) {
			return i, true
		}
	}

	return -1, false
}

func Contains[T comparable](list []T, val T) bool {
	_, ok := FindIndex(list, func(row T) bool { return row == val })

	return ok
}

func ContainsSome[v comparable](list []v, vals []v) bool {
	for _, val := range vals {
		if Contains(list, val) {
			return true
		}
	}

	return false
}

func Shuffle[T any](list []T) []T {
	size := len(list)

	if size < 2 {
		return list
	}

	max := int32(size - 1)

	for i := range list {
		j, _ := RandInt(0, max)
		list[i], list[j] = list[j], list[i]
	}

	return list
}

func ToLower(list []string) []string {
	result := make([]string, len(list))

	for i, item := range list {
		result[i] = strings.ToLower(item)
	}

	return result
}
