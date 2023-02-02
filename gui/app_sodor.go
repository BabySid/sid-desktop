package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/BabySid/gobase"
	"sid-desktop/common"
	"sid-desktop/storage"
	"sid-desktop/theme"
)

type sodorInterface interface {
	CreateView() fyne.CanvasObject
}

var _ sodorInterface = (*sodorAdapter)(nil)

type sodorAdapter struct {
	content fyne.CanvasObject
}

func (d sodorAdapter) CreateView() fyne.CanvasObject {
	panic(any("implement CreateView"))
}

var _ appInterface = (*appSodor)(nil)

type appSodor struct {
	appAdapter
}

func (as *appSodor) LazyInit() error {
	err := storage.GetAppSodorDB().Open(globalWin.app.Storage().RootURI().Path())
	if err != nil {
		return err
	}
	gobase.RegisterAtExit(storage.GetAppSodorDB().Close)

	as.initDB()

	tabs := container.NewAppTabs()

	registerAppTabs(tabs, theme.AppSodorJobTabName, theme.ResourceJobsIcon, &sodorJobs{})
	registerAppTabs(tabs, theme.AppSodorThomasTabName, theme.ResourceTrainIcon, &sodorThomas{})
	registerAppTabs(tabs, theme.AppSodorAlertGroupTabName, theme.ResourceNoticeIcon, &sodorAlertGroup{})
	registerAppTabs(tabs, theme.AppSodorAlertPluginTabName, theme.ResourceAlertIcon, &sodorAlertPlugin{})
	registerAppTabs(tabs, theme.AppSodorFatCtrlTabName, theme.ResourceFatCtrlIcon, &sodorFatController{})

	tabs.SetTabLocation(container.TabLocationLeading)

	tabs.OnSelected = func(item *container.TabItem) {
		// TODO cannot reappear
		// Avoid docTabs invalidation due to theme switching
		item.Content.Refresh()
	}
	as.tabItem = container.NewTabItemWithIcon(theme.AppSodorName, theme.ResourceSodorIcon, tabs)

	if common.GetSodorClient().GetFatCrl().ID == 0 {
		dialog.ShowInformation("", theme.AppSodorInitFatCtlAddrMessage, globalWin.win)
		setFatCtrlSelected(tabs)
	}

	return nil
}

func (as *appSodor) GetAppName() string {
	return theme.AppSodorName
}

func (as *appSodor) ShortCut() fyne.Shortcut {
	return &desktop.CustomShortcut{KeyName: fyne.Key5, Modifier: fyne.KeyModifierAlt}
}

func (as *appSodor) Icon() fyne.Resource {
	return theme.ResourceSodorIcon
}

func (as *appSodor) initDB() {
	need, err := storage.GetAppSodorDB().NeedInit()
	if err != nil {
		printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
		return
	}

	if need {
		err = storage.GetAppSodorDB().Init()
		if err != nil {
			printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
			return
		}
	}

	ctrl, err := storage.GetAppSodorDB().LoadFatCtl()
	if err != nil {
		printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
	} else if ctrl != nil {
		if err = common.GetSodorClient().SetFatCtrlAddr(*ctrl); err != nil {
			printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
		}
	}
}

func registerAppTabs(tabs *container.AppTabs, name string, icon fyne.Resource, s sodorInterface) {
	tabs.Append(container.NewTabItemWithIcon(name, icon, s.CreateView()))
}

func setFatCtrlSelected(tabs *container.AppTabs) {
	for _, item := range tabs.Items {
		if item.Text == theme.AppSodorFatCtrlTabName {
			tabs.Select(item)
			return
		}
	}
}
