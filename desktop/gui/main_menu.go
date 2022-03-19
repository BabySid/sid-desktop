package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	sidTheme "sid-desktop/desktop/theme"
)

type mainMenu struct {
	*fyne.MainMenu

	sysMenu  *fyne.Menu
	appMenus []*fyne.MenuItem

	quit *fyne.MenuItem

	optMenu    *fyne.Menu
	themeOpt   *fyne.MenuItem
	themeDark  *fyne.MenuItem
	themeLight *fyne.MenuItem
	fullScreen *fyne.MenuItem

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
				printErr(fmt.Errorf(sidTheme.RunAppFailedFormat, app.GetAppName(), err))
			}
		})
	}

	// Avoid fyne to add built-in Quit menu
	mm.quit = fyne.NewMenuItem(sidTheme.MenuSysQuit, globalWin.quitHandle)
	mm.quit.IsQuit = true
	mm.sysMenu = fyne.NewMenu(sidTheme.MenuSys, mm.appMenus...)
	mm.sysMenu.Items = append(mm.sysMenu.Items, fyne.NewMenuItemSeparator(), mm.quit)

	// Option-Theme
	mm.themeDark = fyne.NewMenuItem(sidTheme.MenuOptThemeDark, func() {
		globalWin.app.Settings().SetTheme(sidTheme.DarkTheme{})
		mm.themeDark.Checked = true
		mm.themeLight.Checked = false
		globalConfig.Theme.Set("__DARK__")
	})
	mm.themeDark.Checked = true
	mm.themeLight = fyne.NewMenuItem(sidTheme.MenuOptThemeLight, func() {
		globalWin.app.Settings().SetTheme(sidTheme.LightTheme{})
		mm.themeDark.Checked = false
		mm.themeLight.Checked = true
		globalConfig.Theme.Set("__LIGHT__")
	})
	mm.themeLight.Checked = true
	mm.themeOpt = fyne.NewMenuItem(sidTheme.MenuOptTheme, nil)
	mm.themeOpt.ChildMenu = fyne.NewMenu("",
		mm.themeDark,
		mm.themeLight,
	)

	// Option-FullScreen
	mm.fullScreen = fyne.NewMenuItem(sidTheme.MenuOptFullScreen, nil)
	mm.fullScreen.Action = func() {
		if globalWin.win.FullScreen() {
			//mm.fullItem.Label = "FullScreen"
			globalWin.win.SetFullScreen(false)
		} else {
			//mm.fullItem.Label = "QuitFullScreen"
			globalWin.win.SetFullScreen(true)
		}
	}

	// Option
	mm.optMenu = fyne.NewMenu(sidTheme.MenuOption,
		mm.themeOpt,
		fyne.NewMenuItemSeparator(),
		mm.fullScreen,
	)

	// Help
	mm.sysLog = fyne.NewMenuItem(sidTheme.MenuHelpLog, func() {
		newLogViewer().Win.Show()
	})
	mm.aboutSelf = fyne.NewMenuItem(sidTheme.MenuHelpAbout, func() {
		showAbout()
	})
	mm.helpMenu = fyne.NewMenu(sidTheme.MenuHelp,
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

// current version of fyne don't support menu refresh
func (mm *mainMenu) resetMenuStatAfterMainWindowShow() {
	if globalWin.app.Preferences().String("theme") == "__DARK__" {
		mm.themeLight.Checked = false
	} else {
		mm.themeDark.Checked = false
	}
}
