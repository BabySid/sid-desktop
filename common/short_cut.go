package common

import (
	"fyne.io/fyne/v2"
	"strings"
)

func ShortCutName(name fyne.Shortcut) string {
	s := name.ShortcutName()
	arr := strings.Split(s, ":")
	if len(arr) == 2 {
		return arr[1]
	}

	return s
}
