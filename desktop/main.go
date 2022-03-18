package main

import (
	"sid-desktop/desktop/gui"
)

// -ldflags -H=windowsgui
func main() {
	mw := gui.NewMainWin()
	mw.Run()
}
