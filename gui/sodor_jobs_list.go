package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	"github.com/BabySid/proto/sodor"
	"sid-desktop/common"
	"sid-desktop/theme"
)

type sodorJobList struct {
	tabItem *container.TabItem

	searchEntry *widget.Entry
	refresh     *widget.Button
	newJob      *widget.Button

	jobHeader      fyne.CanvasObject
	jobContentList *widget.List

	jobListBinding binding.UntypedList
	jobListCache   *common.JobsWrapper

	createJobHandle    func()
	editJobHandle      func(jobID int32)
	viewInstanceHandle func(job *sodor.Job)
}

func newSodorJobList() *sodorJobList {
	s := sodorJobList{}

	s.searchEntry = widget.NewEntry()
	s.searchEntry.SetPlaceHolder(theme.AppSodorJobSearchText)
	s.searchEntry.OnChanged = s.searchJobs

	s.refresh = widget.NewButton(theme.AppPageRefresh, func() {
		s.loadJobList()
	})
	s.newJob = widget.NewButtonWithIcon(theme.AppSodorCreateJob, theme.ResourceAddIcon, func() {
		if s.createJobHandle != nil {
			s.createJobHandle()
		}
	})

	s.jobListBinding = binding.NewUntypedList()
	s.createJobList()

	s.tabItem = container.NewTabItemWithIcon(theme.AppSodorJobListName, theme.ResourceJobsIcon, nil)
	s.tabItem.Content = container.NewBorder(
		container.NewGridWithColumns(2, s.searchEntry, container.NewHBox(layout.NewSpacer(), s.refresh, s.newJob)),
		nil, nil, nil,
		container.NewBorder(s.jobHeader, nil, nil, nil, s.jobContentList))

	go s.loadJobList()
	return &s
}

func (s *sodorJobList) createJobList() {
	s.jobHeader = container.NewBorder(nil, nil,
		widget.NewLabelWithStyle(theme.AppSodorJobListHeader1, fyne.TextAlignLeading, fyne.TextStyle{}),
		nil,
		container.NewGridWithColumns(4,
			widget.NewLabelWithStyle(theme.AppSodorJobListHeader2, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(theme.AppSodorJobListHeader3, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(theme.AppSodorJobListHeader4, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{})),
	)

	s.jobContentList = widget.NewListWithData(
		s.jobListBinding,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil,
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				nil,
				container.NewGridWithColumns(4,
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					container.NewHBox(
						layout.NewSpacer(),
						widget.NewButtonWithIcon(theme.AppSodorJobListOp1, theme.ResourceEditIcon, nil),
						widget.NewButtonWithIcon(theme.AppSodorJobListOp2, theme.ResourceInstanceIcon, nil),
						widget.NewButtonWithIcon(theme.AppSodorJobListOp3, theme.ResourceRmIcon, nil),
					)),
			)
		},
		func(data binding.DataItem, item fyne.CanvasObject) {
			o, _ := data.(binding.Untyped).Get()
			job := o.(*sodor.Job)

			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%d", job.Id))
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(job.Name)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(gobase.FormatTimeStamp(int64(job.CreateAt)))
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText(gobase.FormatTimeStamp(int64(job.UpdateAt)))

			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*fyne.Container).Objects[1].(*widget.Button).SetText(theme.AppSodorJobListOp1)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
				if s.editJobHandle != nil {
					s.editJobHandle(job.Id)
				}
			}
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*fyne.Container).Objects[2].(*widget.Button).SetText(theme.AppSodorJobListOp2)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*fyne.Container).Objects[2].(*widget.Button).OnTapped = func() {
				if s.viewInstanceHandle != nil {
					s.viewInstanceHandle(job)
				}
			}

			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*fyne.Container).Objects[3].(*widget.Button).SetText(theme.AppSodorJobListOp3)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*fyne.Container).Objects[3].(*widget.Button).OnTapped = func() {
				req := sodor.Job{
					Id: job.Id,
				}
				resp := sodor.JobReply{}
				if err := common.GetSodorClient().Call(common.DeleteJob, &req, &resp); err != nil {
					printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
					return
				}
				s.loadJobList()
			}
		},
	)
}

func (s *sodorJobList) searchJobs(name string) {
	if name == "" {
		if s.jobListCache != nil {
			s.jobListBinding.Set(s.jobListCache.AsInterfaceArray())
		}
	} else {
		if s.jobListCache != nil {
			rs := s.jobListCache.Find(name)
			s.jobListBinding.Set(rs.AsInterfaceArray())
		}
	}
}

func (s *sodorJobList) loadJobList() {
	err := common.GetSodorCache().LoadJobs()
	if err != nil {
		printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
		return
	}
	s.jobListCache = common.NewJobsWrapper(common.GetSodorCache().GetJobs())
	s.jobListBinding.Set(s.jobListCache.AsInterfaceArray())
}
