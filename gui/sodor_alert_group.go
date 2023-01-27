package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"github.com/BabySid/proto/sodor"
	"sid-desktop/theme"
)

var _ sodorInterface = (*sodorAlertGroup)(nil)

type sodorAlertGroup struct {
	sodorAdapter

	docs *container.DocTabs

	alertGroupList *sodorAlertGroupList
}

func (s *sodorAlertGroup) CreateView() fyne.CanvasObject {
	if s.content != nil {
		return s.content
	}

	s.alertGroupList = newSodorAlertGroupList()
	s.alertGroupList.viewHistoryHandle = s.viewAlertPluginInstanceHistory
	s.docs = container.NewDocTabs()
	s.docs.Append(s.alertGroupList.GetTabItem())
	s.docs.SetTabLocation(container.TabLocationTop)
	s.docs.CloseIntercept = func(item *container.TabItem) {
		if item.Text != theme.AppSodorAlertGroupTabName {
			s.docs.Remove(item)
		} else {
			dialog.ShowInformation(theme.CannotCloseTitle, theme.AppSodorAlertGroupListCannotCloseMsg, globalWin.win)
		}
	}

	s.docs.OnSelected = func(item *container.TabItem) {
		// TODO cannot reappear
		// Avoid docTabs invalidation due to theme switching
		item.Content.Refresh()
	}

	s.content = s.docs
	return s.content
}

func (s *sodorAlertGroup) viewAlertPluginInstanceHistory(group *sodor.AlertGroup) {
	info := newSodorAlertGroupHistory(group)
	s.docs.Append(info.tabItem)
	s.docs.Select(info.tabItem)
}
