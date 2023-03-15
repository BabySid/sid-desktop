package common

func EllipsisString(str string, max int) string {
	if len(str) >= max {
		return str[:max] + "..."
	}
	return str
}
