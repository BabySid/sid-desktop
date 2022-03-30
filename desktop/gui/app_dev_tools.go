package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	sidTheme "sid-desktop/desktop/theme"
)

var (
	devToolIndex = map[string][]string{
		"":                               {sidTheme.AppDevToolsTextProcName, sidTheme.AppDevToolsCliName, sidTheme.AppDevToolsDateTimeName},
		sidTheme.AppDevToolsTextProcName: {sidTheme.AppDevToolsJsonProcName, sidTheme.AppDevToolsBase64ProcName},
		sidTheme.AppDevToolsCliName:      {sidTheme.AppDevToolsHttpCliName},
	}

	devTools = map[string]devToolInterface{
		sidTheme.AppDevToolsJsonProcName:   &devToolJson{},
		sidTheme.AppDevToolsBase64ProcName: &devToolBase64{},
		sidTheme.AppDevToolsDateTimeName:   &devToolDateTime{},
	}
)

type devToolInterface interface {
	CreateView() fyne.CanvasObject
}

var _ appInterface = (*appDevTools)(nil)

type appDevTools struct {
	contTree *widget.Tree
	content  *fyne.Container
	tabItem  *container.TabItem
}

func (adt *appDevTools) LazyInit() error {
	adt.contTree = widget.NewTree(
		func(id widget.TreeNodeID) []widget.TreeNodeID {
			return devToolIndex[id]
		}, func(id widget.TreeNodeID) bool {
			children, ok := devToolIndex[id]
			return ok && len(children) > 0
		}, func(b bool) fyne.CanvasObject {
			return widget.NewLabel("")
		}, func(id widget.TreeNodeID, b bool, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(id)
		})

	adt.contTree.OnSelected = func(id widget.TreeNodeID) {
		if obj, ok := devTools[id]; ok {
			adt.content.Objects = []fyne.CanvasObject{obj.CreateView()}
			adt.content.Refresh()
		}
	}

	adt.tabItem = container.NewTabItemWithIcon(sidTheme.AppDevToolsName, sidTheme.ResourceDevToolsIcon, nil)
	adt.content = container.NewMax()
	panel := container.NewHSplit(adt.contTree, adt.content)
	panel.SetOffset(0.3)
	adt.tabItem.Content = panel

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
