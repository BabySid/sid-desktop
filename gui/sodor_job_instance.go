package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"sid-desktop/theme"
)

type sodorJobInstance struct {
	tid int32

	tabItem *container.TabItem

	instanceCard        *widget.Card
	instanceListBinding binding.UntypedList
	instanceHeader      *widget.List
	instanceList        *widget.List

	viewTaskInstanceHandle func(int32)
}

func newSodorJobInstance(id int32) *sodorJobInstance {
	ins := sodorJobInstance{}
	ins.tid = id
	ins.buildJobInstanceInfo()

	ins.tabItem = container.NewTabItem(theme.AppSodorJobInfoJobInstance, nil)
	ins.tabItem.Content = container.NewBorder(
		nil, nil, nil, nil,
		ins.instanceCard)
	return &ins
}

func (s *sodorJobInstance) buildJobInstanceInfo() {
	s.instanceListBinding = binding.NewUntypedList()
	s.instanceListBinding.Set([]interface{}{1, 2, 3})

	s.instanceHeader = widget.NewList(
		func() int {
			return 1
		},
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil,
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				nil,
				container.NewGridWithColumns(6,
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
				),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(theme.AppSodorJobInfoJobInstanceID)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(theme.AppSodorJobInfoJobInstanceScheduleTime)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(theme.AppSodorJobInfoJobInstanceStartTime)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText(theme.AppSodorJobInfoJobInstanceStopTime)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*widget.Label).SetText(theme.AppSodorJobInfoJobInstanceExitCode)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*widget.Label).SetText(theme.AppSodorJobInfoJobInstanceExitMsg)
		},
	)

	s.instanceList = widget.NewListWithData(
		s.instanceListBinding,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil,
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				nil,
				container.NewGridWithColumns(6,
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					container.NewHBox(
						layout.NewSpacer(),
						widget.NewButtonWithIcon(theme.AppSodorJobListOp2, theme.ResourceInstanceIcon, nil),
					),
				),
			)
		},
		func(data binding.DataItem, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText("id")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText("2022-12-12 20:20:33")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText("2022-12-12 20:20:34")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText("2022-12-12 20:20:35")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*widget.Label).SetText("0")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*widget.Label).SetText("{\"output\":\"hello world\"}")

			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[5].(*fyne.Container).Objects[1].(*widget.Button).SetText(theme.AppSodorJobListOp2)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[5].(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
				if s.viewTaskInstanceHandle != nil {
					s.viewTaskInstanceHandle(12345)
				}
			}
		},
	)

	s.instanceCard = widget.NewCard("", theme.AppSodorJobInfoJobInstance,
		container.NewScroll(container.NewBorder(s.instanceHeader, nil, nil, nil, s.instanceList)))
}
