package gui

import (
	"fyne.io/fyne/v2"
	"sid-desktop/theme"
)

type sysTrayMenu struct {
	trayMeny *fyne.Menu
	showMenu *fyne.MenuItem
}

func newSysTrayMenu() *sysTrayMenu {
	sm := sysTrayMenu{}

	sm.showMenu = fyne.NewMenuItem(theme.SysTrayMenuHide, func() {
		if globalWin.wStat.shown {
			globalWin.hideWin()
		} else {
			globalWin.showWin()
		}
	})

	sm.trayMeny = fyne.NewMenu(theme.MenuSys, sm.showMenu)
	return &sm
}

func (s *sysTrayMenu) refreshMenu() {
	if globalWin.wStat.shown {
		s.showMenu.Label = theme.SysTrayMenuHide
	} else {
		s.showMenu.Label = theme.SysTrayMenuShow
	}
	s.trayMeny.Refresh()
}
