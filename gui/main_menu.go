package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"sid-desktop/theme"
)

type mainMenu struct {
	*fyne.MainMenu

	sysMenu  *fyne.Menu
	appMenus []*fyne.MenuItem

	quit *fyne.MenuItem

	optMenu      *fyne.Menu
	themeOpt     *fyne.MenuItem
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

	// System
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

	// Option-Theme
	mm.themeDark = fyne.NewMenuItem(theme.MenuOptThemeDark, func() {
		globalWin.app.Settings().SetTheme(theme.DarkTheme)
		mm.themeDark.Checked = true
		mm.themeLight.Checked = false
		_ = globalConfig.Theme.Set("__DARK__")
		mm.Refresh()
	})
	mm.themeDark.Checked = true
	mm.themeLight = fyne.NewMenuItem(theme.MenuOptThemeLight, func() {
		globalWin.app.Settings().SetTheme(theme.LightTheme)
		mm.themeDark.Checked = false
		mm.themeLight.Checked = true
		_ = globalConfig.Theme.Set("__LIGHT__")
		mm.Refresh()
	})
	mm.themeLight.Checked = true
	t, _ := globalConfig.Theme.Get()
	if t == "__DARK__" {
		mm.themeLight.Checked = false
	} else {
		mm.themeDark.Checked = false
	}

	mm.themeOpt = fyne.NewMenuItem(theme.MenuOptTheme, nil)
	mm.themeOpt.ChildMenu = fyne.NewMenu("",
		mm.themeDark,
		mm.themeLight,
	)

	// Option-FullScreen
	mm.fullScreen = fyne.NewMenuItem(theme.MenuOptFullScreen, nil)
	mm.fullScreen.Action = func() {
		if globalWin.win.FullScreen() {
			mm.fullScreen.Label = "FullScreen"
			globalWin.win.SetFullScreen(false)
		} else {
			mm.fullScreen.Label = "QuitFullScreen"
			globalWin.win.SetFullScreen(true)
		}
		mm.Refresh()
	}
	if globalWin.win.FullScreen() {
		mm.fullScreen.Label = "QuitFullScreen"
	} else {
		mm.fullScreen.Label = "FullScreen"
	}

	// Option-HideWhenQuit
	mm.hideWhenQuit = fyne.NewMenuItem(theme.MenuOptHideWhenQuit, nil)
	mm.hideWhenQuit.Checked = true
	mm.hideWhenQuit.Action = func() {
		hide, _ := globalConfig.HideWhenQuit.Get()
		if hide {
			mm.hideWhenQuit.Checked = false
			_ = globalConfig.HideWhenQuit.Set(false)
		} else {
			mm.hideWhenQuit.Checked = true
			_ = globalConfig.HideWhenQuit.Set(true)
		}
		mm.Refresh()
	}
	hide, _ := globalConfig.HideWhenQuit.Get()
	if !hide {
		mm.hideWhenQuit.Checked = false
	}

	// Option
	mm.optMenu = fyne.NewMenu(theme.MenuOption,
		mm.themeOpt,
		fyne.NewMenuItemSeparator(),
		mm.fullScreen,
		fyne.NewMenuItemSeparator(),
		mm.hideWhenQuit,
	)

	// Help
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

	mm.MainMenu = fyne.NewMainMenu(
		mm.sysMenu,
		mm.optMenu,
		mm.helpMenu,
	)

	return &mm
}
