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
		&appDevTools{},
	}
)

var _ appInterface = (*appAdapter)(nil)

type appAdapter struct {
	tabItem *container.TabItem
}

func (a appAdapter) LazyInit() error {
	panic("implement LazyInit")
}

func (a appAdapter) GetTabItem() *container.TabItem {
	return a.tabItem
}

func (a appAdapter) GetAppName() string {
	panic("implement GetAppName()")
}

func (a appAdapter) OpenDefault() bool {
	return false
}

func (a appAdapter) OnClose() bool {
	return true
}
