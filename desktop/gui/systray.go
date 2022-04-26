package gui

import (
	"fyne.io/systray"
	sidTheme "sid-desktop/desktop/theme"
)

type sysTray struct {
	showMenu *systray.MenuItem
	quitMenu *systray.MenuItem
}

func newSysTray() *sysTray {
	return &sysTray{showMenu: nil}
}

func (s *sysTray) run() {
	go func() {
		systray.Run(s.systrayReady, s.systrayExit)
	}()
}

func (s *sysTray) setShowMenu() {
	s.showMenu.SetTitle(sidTheme.SysTrayMenuShow)
	s.showMenu.SetTooltip(sidTheme.SysTrayMenuShowTooltip)
}

func (s *sysTray) setHideMenu() {
	s.showMenu.SetTitle(sidTheme.SysTrayMenuHide)
	s.showMenu.SetTooltip(sidTheme.SysTrayMenuHideTooltip)
}

func (s *sysTray) systrayReady() {
	systray.SetIcon(sidTheme.ResourceSystrayIcon.Content())
	systray.SetTitle(sidTheme.AppTitle)
	systray.SetTooltip(sidTheme.AppTitle)

	s.showMenu = systray.AddMenuItem(sidTheme.SysTrayMenuHide, sidTheme.SysTrayMenuHideTooltip)
	systray.AddSeparator()
	s.quitMenu = systray.AddMenuItem(sidTheme.MenuSysQuit, sidTheme.SysTrayMenuQuitTooltip)

	go func() {
		for {
			select {
			case <-s.showMenu.ClickedCh:
				if globalWin.wStat.shown {
					globalWin.hideWin()
				} else {
					globalWin.showWin()
				}
			case <-s.quitMenu.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func (s *sysTray) systrayExit() {
	globalWin.closeWin()
}
