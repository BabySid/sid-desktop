package gui

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2/dialog"
	"github.com/BabySid/gobase"
	"log"
	sidTheme "sid-desktop/theme"
)

func printErr(err error) {
	log.Printf(sidTheme.InternalErrorFormat, err)

	str := err.Error()

	target := ""
	const maxLen = 64

	n := len(str) / maxLen
	m := len(str) % maxLen
	for i := 0; i < n; i++ {
		begin := i * maxLen
		end := begin + maxLen
		if i > 0 {
			target += "\n"
		}
		target += str[begin:end]
	}
	if m > 0 {
		if n > 0 {
			target += "\n"
		}
		target += str[n*maxLen : n*maxLen+m]
	}

	dialog.ShowError(errors.New(target), globalWin.win)
}

func assertErr(cond bool, a ...interface{}) {
	if !cond {
		printErr(fmt.Errorf(fmt.Sprint(a)))
		gobase.AssertHere()
	}
}
