package gui

import (
	"fyne.io/fyne/v2/container"
	sidTheme "sid-desktop/desktop/theme"
)

var _ appInterface = (*appDevTools)(nil)

type appDevTools struct {
	tabItem *container.TabItem
}

func (adt *appDevTools) LazyInit() error {
	return nil
}

func (adt *appDevTools) GetTabItem() *container.TabItem {
	return adt.tabItem
}

func (adt *appDevTools) GetAppName() string {
	return sidTheme.AppDevToolsName
}

func (adt *appDevTools) OpenDefault() bool {
	return false
}

func (adt *appDevTools) OnClose() bool {
	return true
}
