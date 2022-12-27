package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
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
	panic("implement CreateView")
}

var _ appInterface = (*appSodor)(nil)

type appSodor struct {
	appAdapter
}

func (as *appSodor) LazyInit() error {
	tabs := container.NewAppTabs()

	registerAppTabs(tabs, theme.AppSodorJobTabName, theme.ResourceJobsIcon, &sodorJobs{})
	registerAppTabs(tabs, theme.AppSodorThomsTabName, theme.ResourceTrainIcon, &sodorThomas{})
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

func registerAppTabs(tabs *container.AppTabs, name string, icon fyne.Resource, s sodorInterface) {
	tabs.Append(container.NewTabItemWithIcon(name, icon, s.CreateView()))
}
