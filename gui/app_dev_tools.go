package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	"sid-desktop/storage"
	"sid-desktop/theme"
)

var (
	devToolIndex = map[string][]string{
		"":                            {theme.AppDevToolsTextProcName, theme.AppDevToolsCliName, theme.AppDevToolsDateTimeName},
		theme.AppDevToolsTextProcName: {theme.AppDevToolsJsonProcName, theme.AppDevToolsBase64ProcName},
		theme.AppDevToolsCliName:      {theme.AppDevToolsHttpCliName},
	}

	devTools = map[string]devToolInterface{
		theme.AppDevToolsJsonProcName:   &devToolJson{},
		theme.AppDevToolsBase64ProcName: &devToolBase64{},
		theme.AppDevToolsDateTimeName:   &devToolDateTime{},
		theme.AppDevToolsHttpCliName:    &devToolHttpClient{},
	}
)

type devToolInterface interface {
	CreateView() fyne.CanvasObject
}

var _ devToolInterface = (*devToolAdapter)(nil)

type devToolAdapter struct {
	content fyne.CanvasObject
}

func (d devToolAdapter) CreateView() fyne.CanvasObject {
	panic("implement CreateView")
}

var _ appInterface = (*appDevTools)(nil)

type appDevTools struct {
	appAdapter
	contTree *widget.Tree
	content  *fyne.Container
}

func (adt *appDevTools) LazyInit() error {
	err := storage.GetAppDevToolDB().Open(globalWin.app.Storage().RootURI().Path())
	if err != nil {
		return err
	}
	gobase.RegisterAtExit(storage.GetAppDevToolDB().Close)

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

	adt.contTree.OpenAllBranches()

	adt.tabItem = container.NewTabItemWithIcon(theme.AppDevToolsName, theme.ResourceDevToolsIcon, nil)
	adt.content = container.NewMax()
	panel := container.NewHSplit(adt.contTree, adt.content)
	panel.SetOffset(0.15)
	adt.tabItem.Content = panel

	go adt.initDB()

	return nil
}

func (adt *appDevTools) GetAppName() string {
	return theme.AppDevToolsName
}

func (adt *appDevTools) initDB() {
	need, err := storage.GetAppDevToolDB().NeedInit()
	if err != nil {
		printErr(fmt.Errorf(theme.AppDevToolsFailedFormat, err))
		return
	}

	if need {
		err = storage.GetAppDevToolDB().Init()
		if err != nil {
			printErr(fmt.Errorf(theme.AppDevToolsFailedFormat, err))
			return
		}
	}
}

func (adt *appDevTools) ShortCut() fyne.Shortcut {
	return &desktop.CustomShortcut{KeyName: fyne.Key5, Modifier: desktop.AltModifier}
}

func (adt *appDevTools) Icon() fyne.Resource {
	return theme.ResourceDevToolsIcon
}
