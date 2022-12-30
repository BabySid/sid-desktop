package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"sid-desktop/theme"
)

type sodorJobList struct {
	tabItem *container.TabItem

	searchEntry *widget.Entry
	newJob      *widget.Button

	jobHeader      *widget.List
	jobContentList *widget.List

	jobListBinding binding.UntypedList

	createJobHandle    func()
	editJobHandle      func()
	viewInstanceHandle func()
}

func newSodorJobList() *sodorJobList {
	s := sodorJobList{}

	s.searchEntry = widget.NewEntry()
	s.searchEntry.SetPlaceHolder(theme.AppSodorJobSearchText)
	s.searchEntry.OnChanged = s.searchJobs

	s.newJob = widget.NewButtonWithIcon(theme.AppSodorCreateJob, theme.ResourceAddIcon, func() {
		if s.createJobHandle != nil {
			s.createJobHandle()
		}
	})

	s.jobListBinding = binding.NewUntypedList()
	s.jobListBinding.Set([]interface{}{1, 2, 3})
	s.createJobList()

	s.tabItem = container.NewTabItemWithIcon(theme.AppSodorJobListName, theme.ResourceJobsIcon, nil)
	s.tabItem.Content = container.NewBorder(
		container.NewGridWithColumns(2, s.searchEntry, container.NewHBox(layout.NewSpacer(), s.newJob)),
		nil, nil, nil,
		container.NewScroll(container.NewBorder(s.jobHeader, nil, nil, nil, s.jobContentList)))
	return &s
}

func (s *sodorJobList) GetText() string {
	return s.tabItem.Text
}

func (s *sodorJobList) GetTabItem() *container.TabItem {
	return s.tabItem
}

func (s *sodorJobList) createJobList() {
	s.jobHeader = widget.NewList(
		func() int {
			return 1
		},
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil,
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				nil,
				container.NewGridWithColumns(4,
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{})),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(theme.AppSodorJobListHeader1)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(theme.AppSodorJobListHeader2)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(theme.AppSodorJobListHeader3)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText(theme.AppSodorJobListHeader4)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*widget.Label).SetText("")
		},
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
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText("1")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText("job name")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText("2022-12-12 23:22:15")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText("2022-12-15 23:22:15")

			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*fyne.Container).Objects[1].(*widget.Button).SetText(theme.AppSodorJobListOp1)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
				if s.editJobHandle != nil {
					s.editJobHandle()
				}
			}
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*fyne.Container).Objects[2].(*widget.Button).SetText(theme.AppSodorJobListOp2)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*fyne.Container).Objects[2].(*widget.Button).OnTapped = func() {
				if s.viewInstanceHandle != nil {
					s.viewInstanceHandle()
				}
			}

			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*fyne.Container).Objects[3].(*widget.Button).SetText(theme.AppSodorJobListOp3)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*fyne.Container).Objects[3].(*widget.Button).OnTapped = func() {
				// rm
			}
		},
	)
}

func (s *sodorJobList) searchJobs(name string) {

}
