package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"sid-desktop/theme"
)

type sodorThomasInstance struct {
	tid int32

	tabItem *container.TabItem

	thomasCard          *widget.Card
	thomasID            *widget.Label
	thomasName          *widget.Label
	thomasVersion       *widget.Label
	thomasProto         *widget.Label
	thomasHost          *widget.Label
	thomasPort          *widget.Label
	thomasPID           *widget.Label
	thomasStartTime     *widget.Label
	thomasHeartBeatTime *widget.Label
	thomasType          *widget.Label
	thomasStatus        *widget.Label

	instanceCard        *widget.Card
	instanceListBinding binding.UntypedList
	instanceHeader      *widget.List
	instanceList        *widget.List
}

func newSodorThomasInstance(id int32) *sodorThomasInstance {
	ins := sodorThomasInstance{}
	ins.tid = id

	ins.buildThomasInfo()
	ins.buildThomasInstanceInfo()

	ins.tabItem = container.NewTabItem(theme.AppSodorThomasInfo, nil)
	ins.tabItem.Content = container.NewBorder(
		ins.thomasCard, nil, nil, nil,
		ins.instanceCard)
	return &ins
}

func (s *sodorThomasInstance) buildThomasInfo() {
	s.thomasID = widget.NewLabel(fmt.Sprintf("%d", s.tid))
	s.thomasName = widget.NewLabel("thomasName")
	s.thomasVersion = widget.NewLabel("thomasVersion")
	s.thomasProto = widget.NewLabel("thomasProto")
	s.thomasHost = widget.NewLabel("thomasHost")
	s.thomasPort = widget.NewLabel("thomasPort")
	s.thomasPID = widget.NewLabel("thomasPID")
	s.thomasStartTime = widget.NewLabel("thomasStartTime")
	s.thomasHeartBeatTime = widget.NewLabel("thomasHeartBeatTime")
	s.thomasType = widget.NewLabel("thomasType")
	s.thomasStatus = widget.NewLabel("thomasStatus")

	infoBox := container.NewVBox()
	infoBox.Add(
		container.NewHBox(
			widget.NewForm(widget.NewFormItem(theme.AppSodorThomasInfoID, s.thomasID)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorThomasInfoName, s.thomasName)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorThomasInfoVersion, s.thomasVersion)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorThomasInfoType, s.thomasType)),
		),
	)
	infoBox.Add(
		container.NewHBox(
			widget.NewForm(widget.NewFormItem(theme.AppSodorThomasInfoProto, s.thomasProto)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorThomasInfoHost, s.thomasHost)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorThomasInfoPort, s.thomasPort)),
		),
	)
	infoBox.Add(
		container.NewHBox(
			widget.NewForm(widget.NewFormItem(theme.AppSodorThomasInfoStartTime, s.thomasStartTime)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorThomasInfoHeartBeatTime, s.thomasHeartBeatTime)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorThomasInfoPID, s.thomasPID)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorThomasInfoStatus, s.thomasStatus)),
		),
	)

	s.thomasCard = widget.NewCard("", theme.AppSodorThomasInfo, infoBox)
}

func (s *sodorThomasInstance) buildThomasInstanceInfo() {
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
				container.NewBorder(nil, nil,
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					nil,
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(theme.AppSodorThomasInfoInstanceID)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(theme.AppSodorThomasInfoInstanceCreateTime)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(theme.AppSodorThomasInfoInstanceMetrics)
		},
	)

	s.instanceList = widget.NewListWithData(
		s.instanceListBinding,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil,
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				nil,
				container.NewBorder(nil, nil,
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					nil,
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{})),
			)
		},
		func(data binding.DataItem, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText("id")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText("2022-12-12 20:20:33")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText("{\"cpu\":123,\"mem\":456}")
		},
	)

	s.instanceCard = widget.NewCard("", theme.AppSodorThomasInfoInstance,
		container.NewScroll(container.NewBorder(s.instanceHeader, nil, nil, nil, s.instanceList)))
}
