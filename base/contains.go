package base

func ContainsString(array []string, val string) int {
	for i, item := range array {
		if item == val {
			return i
		}
	}

	return -1
}
