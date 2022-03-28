package gui

import (
	"fmt"
	"fyne.io/fyne/v2/widget"
	sidTheme "sid-desktop/desktop/theme"
)

type toolBar struct {
	toolbar *widget.Toolbar
}

func newToolBar() *toolBar {
	var tb toolBar

	tb.toolbar = widget.NewToolbar(
		widget.NewToolbarAction(sidTheme.ResourceWelIcon, func() {
			openApp(sidTheme.AppWelcomeName)
		}),
		widget.NewToolbarAction(sidTheme.ResourceLauncherIcon, func() {
			openApp(sidTheme.AppLauncherName)
		}),
		widget.NewToolbarAction(sidTheme.ResourceFavoritesIcon, func() {
			openApp(sidTheme.AppFavoritesName)
		}),
		widget.NewToolbarAction(sidTheme.ResourceScriptRunnerIcon, func() {
			openApp(sidTheme.AppScriptRunnerName)
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(sidTheme.ResourceLogViewerIcon, func() {
			newLogViewer().Win.Show()
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(sidTheme.ResourceAboutIcon, func() {
			showAbout()
		}),
	)

	return &tb
}

func openApp(name string) {
	for _, app := range appRegister {
		if app.GetAppName() == name {
			err := globalWin.at.openApp(app)
			if err != nil {
				printErr(fmt.Errorf(sidTheme.RunAppFailedFormat, app.GetAppName(), err))
			}
		}
	}
}
