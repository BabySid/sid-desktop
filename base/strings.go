package base

import (
	"unicode/utf8"
)

func CutUTF8(str string, start int, end int) string {
	if end <= utf8.RuneCountInString(str) {
		aft := ""
		i := 0
		for _, v := range str {
			if i >= start && i < end {
				aft += string(v)
			}
			if i >= end {
				break
			}
			i++
		}

		return aft
	}
	return str
}
