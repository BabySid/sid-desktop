package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	sidTheme "sid-desktop/desktop/theme"
)

var _ appInterface = (*appDevTools)(nil)

type appDevTools struct {
	contTree *widget.Tree
	tabItem  *container.TabItem
}

func (adt *appDevTools) LazyInit() error {
	adt.contTree = widget.NewTree(
		func(id widget.TreeNodeID) []widget.TreeNodeID {
			return []string{"a", "b", "c"}
		}, func(id widget.TreeNodeID) bool {
			return true
		}, func(b bool) fyne.CanvasObject {
			return widget.NewLabel("a")
		}, func(id widget.TreeNodeID, b bool, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText("title")
		})

	adt.tabItem = container.NewTabItemWithIcon(sidTheme.AppDevToolsName, sidTheme.ResourceDevToolsIcon, nil)
	adt.tabItem.Content = container.NewHSplit(adt.contTree, widget.NewLabel("this is content"))
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
