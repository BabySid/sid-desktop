package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"github.com/BabySid/proto/sodor"
	"sid-desktop/theme"
	"sync"
)

var _ sodorInterface = (*sodorThomas)(nil)

type sodorThomas struct {
	sodorAdapter

	docs       *container.DocTabs
	thomasList *sodorThomasList

	contentPages sync.Map
}

func (s *sodorThomas) CreateView() fyne.CanvasObject {
	if s.content != nil {
		return s.content
	}

	s.thomasList = newSodorThomasList()
	s.thomasList.viewInstanceHandle = s.viewThomasInstance
	s.docs = container.NewDocTabs()
	s.docs.Append(s.thomasList.GetTabItem())
	s.docs.SetTabLocation(container.TabLocationTop)
	s.docs.CloseIntercept = func(item *container.TabItem) {
		if item.Text != theme.AppSodorThomasTabName {
			s.docs.Remove(item)
			if page, ok := s.contentPages.LoadAndDelete(item.Text); ok {
				page.(sodorContentPage).OnClose()
			}
		} else {
			dialog.ShowInformation(theme.CannotCloseTitle, theme.AppSodorThomasListCannotCloseMsg, globalWin.win)
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

func (s *sodorThomas) viewThomasInstance(info *sodor.ThomasInfo) {
	for _, item := range s.docs.Items {
		if item.Text == getTabNameOfThomasInstance(info) {
			s.docs.Select(item)
			return
		}
	}

	thomas := newSodorThomasInstance(info)
	s.docs.Append(thomas.tabItem)
	s.docs.Select(thomas.tabItem)

	s.contentPages.Store(thomas.tabItem.Text, thomas)
}
