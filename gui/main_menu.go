package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"sid-desktop/common"
	"sid-desktop/theme"
)

type mainMenu struct {
	*fyne.MainMenu

	sysMenu  *fyne.Menu
	appMenus []*fyne.MenuItem

	quit *fyne.MenuItem

	optMenu      *fyne.Menu
	themeOpt     *fyne.MenuItem
	themeDefault *fyne.MenuItem
	themeDark    *fyne.MenuItem
	themeLight   *fyne.MenuItem
	fullScreen   *fyne.MenuItem
	hideWhenQuit *fyne.MenuItem

	helpMenu  *fyne.Menu
	sysLog    *fyne.MenuItem
	aboutSelf *fyne.MenuItem
}

func newMainMenu() *mainMenu {
	var mm mainMenu
	setSystemMenu(&mm)
	setOptMenu(&mm)
	setHelpMenu(&mm)
	mm.MainMenu = fyne.NewMainMenu(
		mm.sysMenu,
		mm.optMenu,
		mm.helpMenu,
	)

	return &mm
}

func changeTheme(themeName string) {
	fmt.Println("切换主题为:", themeName)
	dialog.ShowConfirm(theme.RestartTitle, theme.RestartMsg, func(b bool) {
		if b {
			_ = common.GetConfig().Theme.Set(themeName)
			globalWin.restart()
		}
	}, globalWin.win)

}

func setThemeMenuCheckStatus(mm *mainMenu) {
	t, _ := common.GetConfig().Theme.Get()
	if t == "__DARK__" {
		mm.themeDark.Checked = true
		mm.themeLight.Checked = false
		mm.themeDefault.Checked = false
	} else if t == "__LIGHT__" {
		mm.themeLight.Checked = true
		mm.themeDark.Checked = false
		mm.themeDefault.Checked = false
	} else if t == "__DEFAULT__" {
		mm.themeDefault.Checked = true
		mm.themeLight.Checked = false
		mm.themeDark.Checked = false
	}
}

func setThemeMenu(mm *mainMenu) {
	// Option-Theme
	mm.themeDefault = fyne.NewMenuItem(theme.MenuOptThemeDefault, func() {
		changeTheme("__DEFAULT__")
	})
	//	mm.themeDefault.Checked = true
	mm.themeDark = fyne.NewMenuItem(theme.MenuOptThemeDark, func() {
		changeTheme("__DARK__")
	})
	//mm.themeDark.Checked = true
	mm.themeLight = fyne.NewMenuItem(theme.MenuOptThemeLight, func() {
		changeTheme("__LIGHT__")
	})
	mm.themeOpt = fyne.NewMenuItem(theme.MenuOptTheme, nil)
	mm.themeOpt.ChildMenu = fyne.NewMenu("",
		mm.themeDefault,
		mm.themeDark,
		mm.themeLight,
	)
	//mm.themeLight.Checked = true
}

func setFullScreenMenu(mm *mainMenu) {
	mm.fullScreen = fyne.NewMenuItem(theme.MenuOptFullScreen, nil)
	mm.fullScreen.Action = func() {
		if globalWin.win.FullScreen() {
			globalWin.win.SetFullScreen(false)
			mm.fullScreen.Label = "全屏"
		} else {
			globalWin.win.SetFullScreen(true)
			mm.fullScreen.Label = "退出全屏"
		}
		mm.Refresh()
	}
	if globalWin.win.FullScreen() {
		mm.fullScreen.Label = "退出全屏"
	} else {
		mm.fullScreen.Label = "全屏"
	}
}

func setHideWhenQuitMenu(mm *mainMenu) {
	mm.hideWhenQuit = fyne.NewMenuItem(theme.MenuOptHideWhenQuit, nil)
	mm.hideWhenQuit.Checked = true
	mm.hideWhenQuit.Action = func() {
		hide, _ := common.GetConfig().HideWhenQuit.Get()
		if hide {
			mm.hideWhenQuit.Checked = false
			_ = common.GetConfig().HideWhenQuit.Set(false)
		} else {
			mm.hideWhenQuit.Checked = true
			_ = common.GetConfig().HideWhenQuit.Set(true)
		}
		mm.Refresh()
	}
	hide, _ := common.GetConfig().HideWhenQuit.Get()
	if !hide {
		mm.hideWhenQuit.Checked = false
	}
}

func setSystemMenu(mm *mainMenu) {
	mm.appMenus = make([]*fyne.MenuItem, len(appRegister))
	for i, app := range appRegister {
		app := app
		mm.appMenus[i] = fyne.NewMenuItem(app.GetAppName(), func() {
			err := globalWin.at.openApp(app)
			if err != nil {
				printErr(fmt.Errorf(theme.RunAppFailedFormat, app.GetAppName(), err))
			}
		})
		mm.appMenus[i].Shortcut = app.ShortCut()
		mm.appMenus[i].Icon = app.Icon()
	}

	// Avoid fyne to add built-in Quit menu
	mm.quit = fyne.NewMenuItem(theme.MenuSysQuit, globalWin.quitHandle)
	mm.quit.IsQuit = true
	mm.sysMenu = fyne.NewMenu(theme.MenuSys, mm.appMenus...)
	mm.sysMenu.Items = append(mm.sysMenu.Items, fyne.NewMenuItemSeparator(), mm.quit)
}

func setHelpMenu(mm *mainMenu) {
	mm.sysLog = fyne.NewMenuItem(theme.MenuHelpLog, func() {
		newLogViewer().Win.Show()
	})
	mm.sysLog.Icon = theme.ResourceLogViewerIcon
	mm.aboutSelf = fyne.NewMenuItem(theme.MenuHelpAbout, func() {
		showAbout()
	})
	mm.aboutSelf.Icon = theme.ResourceAboutIcon
	mm.helpMenu = fyne.NewMenu(theme.MenuHelp,
		mm.sysLog,
		mm.aboutSelf,
	)
}

func setOptMenu(mm *mainMenu) {
	setThemeMenu(mm)
	setThemeMenuCheckStatus(mm)
	setFullScreenMenu(mm)
	setHideWhenQuitMenu(mm)
	mm.optMenu = fyne.NewMenu(theme.MenuOption,
		mm.themeOpt,
		fyne.NewMenuItemSeparator(),
		mm.fullScreen,
		fyne.NewMenuItemSeparator(),
		mm.hideWhenQuit,
	)
}
