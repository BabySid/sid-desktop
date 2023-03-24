package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

type appInterface interface {
	LazyInit() error // run when click on menu or toolbar to run the app
	GetTabItem() *container.TabItem
	GetAppName() string
	OpenDefault() bool
	OnClose() bool
	ShortCut() fyne.Shortcut
	Icon() fyne.Resource
}

var (
	appRegister = []appInterface{
		&appWelcome{},
		&appLauncher{},
		&appFavorites{},
		&appMarkDown{},
		&appDevTools{},
		&appSodor{},
	}
)

var _ appInterface = (*appAdapter)(nil)

type appAdapter struct {
	tabItem *container.TabItem
}

func (a appAdapter) ShortCut() fyne.Shortcut {
	panic("implement ShortCut")
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

func (a appAdapter) Icon() fyne.Resource {
	return theme.FyneLogo()
}
