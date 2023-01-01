package gui

import (
	"fyne.io/systray"
	"sid-desktop/theme"
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
	s.showMenu.SetTitle(theme.SysTrayMenuShow)
	s.showMenu.SetTooltip(theme.SysTrayMenuShowTooltip)
}

func (s *sysTray) setHideMenu() {
	s.showMenu.SetTitle(theme.SysTrayMenuHide)
	s.showMenu.SetTooltip(theme.SysTrayMenuHideTooltip)
}

func (s *sysTray) systrayReady() {
	systray.SetIcon(theme.ResourceSystrayIcon.Content())
	systray.SetTitle(theme.AppTitle)
	systray.SetTooltip(theme.AppTitle)

	s.showMenu = systray.AddMenuItem(theme.SysTrayMenuHide, theme.SysTrayMenuHideTooltip)
	systray.AddSeparator()
	s.quitMenu = systray.AddMenuItem(theme.MenuSysQuit, theme.SysTrayMenuQuitTooltip)

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

func (s *sysTray) Quit() {
	systray.Quit()
}
