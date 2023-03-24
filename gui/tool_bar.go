package gui

import (
	"fmt"
	"fyne.io/fyne/v2/widget"
	"sid-desktop/theme"
)

type toolBar struct {
	toolbar *widget.Toolbar
}

func newToolBar() *toolBar {
	var tb toolBar

	tb.toolbar = widget.NewToolbar(
		widget.NewToolbarAction(theme.ResourceWelIcon, func() {
			openApp(theme.AppWelcomeName)
		}),
		widget.NewToolbarAction(theme.ResourceLauncherIcon, func() {
			openApp(theme.AppLauncherName)
		}),
		widget.NewToolbarAction(theme.ResourceFavoritesIcon, func() {
			openApp(theme.AppFavoritesName)
		}),
		widget.NewToolbarAction(theme.ResourceDevToolsIcon, func() {
			openApp(theme.AppDevToolsName)
		}),
		widget.NewToolbarAction(theme.ResourceSodorIcon, func() {
			openApp(theme.AppSodorName)
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ResourceLogViewerIcon, func() {
			newLogViewer().Win.Show()
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.ResourceAboutIcon, func() {
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
				printErr(fmt.Errorf(theme.RunAppFailedFormat, app.GetAppName(), err))
			}
		}
	}
}
