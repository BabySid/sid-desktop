package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"sid-desktop/theme"
)

var _ sodorInterface = (*sodorThomas)(nil)

type sodorThomas struct {
	sodorAdapter

	docs       *container.DocTabs
	thomasList *sodorThomasList
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
		if item.Text != theme.AppSodorThomsTabName {
			s.docs.Remove(item)
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

func (s *sodorThomas) viewThomasInstance(thomasID int32) {
	thomas := newSodorThomasInstance(thomasID)
	s.docs.Append(thomas.tabItem)
	s.docs.Select(thomas.tabItem)
}
