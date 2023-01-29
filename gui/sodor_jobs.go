package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"github.com/BabySid/proto/sodor"
	"sid-desktop/theme"
)

var _ sodorInterface = (*sodorJobs)(nil)

type sodorJobs struct {
	sodorAdapter

	docs    *container.DocTabs
	jobList *sodorJobList
}

func (s *sodorJobs) CreateView() fyne.CanvasObject {
	if s.content != nil {
		return s.content
	}

	s.jobList = newSodorJobList()
	s.jobList.createJobHandle = s.createJob
	s.jobList.editJobHandle = s.editJob
	s.jobList.viewInstanceHandle = s.viewJobInstance

	s.docs = container.NewDocTabs()
	s.docs.Append(s.jobList.tabItem)
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

func (s *sodorJobs) createJob() {
	info := newSodorJobInfo(0)
	s.docs.Append(info.tabItem)
	s.docs.Select(info.tabItem)
	info.okHandle = func() {
		s.jobList.loadJobList()
		s.docs.Remove(info.tabItem)
	}
	info.dismissHandle = func() {
		s.docs.Remove(info.tabItem)
	}
}

func (s *sodorJobs) editJob(id int32) {
	info := newSodorJobInfo(id)
	s.docs.Append(info.tabItem)
	s.docs.Select(info.tabItem)
	info.okHandle = func() {
		s.docs.Remove(info.tabItem)
	}
	info.dismissHandle = func() {
		s.docs.Remove(info.tabItem)
	}
}

func (s *sodorJobs) viewJobInstance(job *sodor.Job) {
	info := newSodorJobInstance(job)
	info.viewTaskInstanceHandle = s.viewTaskInstance
	s.docs.Append(info.tabItem)
	s.docs.Select(info.tabItem)
}

func (s *sodorJobs) viewTaskInstance(param taskInstanceParam) {
	info := newSodorJobTaskInstance(param)
	s.docs.Append(info.tabItem)
	s.docs.Select(info.tabItem)
}
