package gui

import "fyne.io/fyne/v2/container"

type appInterface interface {
	LazyInit() error // run when click on menu or toolbar to run the app
	GetTabItem() *container.TabItem
	GetAppName() string
	OpenDefault() bool
	OnClose() bool
}

var (
	appRegister = []appInterface{
		&appWelcome{},
		&appLauncher{},
		&appFavorites{},
		&appScriptRunner{},
	}
)
