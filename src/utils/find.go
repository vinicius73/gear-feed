package utils

func FindStr(list []string, val string) (int, bool) {
	for i, item := range list {
		if item == val {
			return i, true
		}
	}

	return -1, false
}
