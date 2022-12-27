package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"sid-desktop/theme"
)

var _ sodorInterface = (*sodorJobs)(nil)

type sodorJobs struct {
	sodorAdapter

	docs     *container.DocTabs
	mainView *sodorJobsMainView
}

func (s *sodorJobs) CreateView() fyne.CanvasObject {
	if s.content != nil {
		return s.content
	}

	s.mainView = newSodorJobsMainView()
	s.docs = container.NewDocTabs()
	s.docs.Append(s.mainView.GetTabItem())
	s.docs.SetTabLocation(container.TabLocationTop)
	s.docs.CloseIntercept = func(item *container.TabItem) {
		if item.Text != theme.AppSodorJobListName {
			s.docs.Remove(item)
		} else {
			dialog.ShowInformation(theme.CannotCloseTitle, theme.AppSodorJobListCannotCloseMsg, globalWin.win)
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
