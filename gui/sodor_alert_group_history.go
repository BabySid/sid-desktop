package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"sid-desktop/theme"
)

type sodorAlertGroupHistory struct {
	grpID int32

	tabItem *container.TabItem

	groupID             *widget.Label
	groupName           *widget.Label
	pluginInstanceNames *widget.CheckGroup
	headerContainer     *fyne.Container

	historyCard        *widget.Card
	historyListBinding binding.UntypedList
	historyHeader      *widget.List
	historyList        *widget.List
}

func newSodorAlertGroupHistory(id int32) *sodorAlertGroupHistory {
	ins := sodorAlertGroupHistory{}
	ins.grpID = id

	ins.buildPluginNames()
	ins.buildHistoryInfo()

	ins.tabItem = container.NewTabItem(theme.AppSodorAddAlertGroupHistory, nil)
	ins.tabItem.Content = container.NewBorder(
		ins.headerContainer, nil, nil, nil,
		ins.historyCard)
	return &ins
}

func (s *sodorAlertGroupHistory) buildPluginNames() {
	s.groupID = widget.NewLabel("12345")
	s.groupName = widget.NewLabel("jobName")

	s.pluginInstanceNames = widget.NewCheckGroup([]string{"plugin1", "plugin2", "plugin3", "plugin4"}, nil)
	s.pluginInstanceNames.Horizontal = true

	s.headerContainer = container.NewVBox(
		container.NewHBox(
			widget.NewForm(widget.NewFormItem(theme.AppSodorCreateAlertGroupID, s.groupID)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorCreateAlertGroupName, s.groupName)),
			layout.NewSpacer()),
	)
	s.headerContainer.Add(widget.NewCard("", theme.AppSodorCreateAlertGroupPlugins, s.pluginInstanceNames))
}

func (s *sodorAlertGroupHistory) buildHistoryInfo() {
	s.historyListBinding = binding.NewUntypedList()
	s.historyListBinding.Set([]interface{}{1, 2, 3})

	s.historyHeader = widget.NewList(
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
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(theme.AppSodorCreateAlertGroupPluginID)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(theme.AppSodorCreateAlertGroupPluginName)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(theme.AppSodorCreateAlertGroupAlertMsg)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText(theme.AppSodorCreateAlertGroupAlertCreateTime)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*widget.Label).SetText(theme.AppSodorCreateAlertGroupStatusMsg)
		},
	)

	s.historyList = widget.NewListWithData(
		s.historyListBinding,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil,
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				nil,
				container.NewGridWithColumns(4,
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{}),
				),
			)
		},
		func(data binding.DataItem, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText("id")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText("group name")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText("alert to xxxxxxxxxxxx")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText("2022-12-12 20:20:33")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*widget.Label).SetText("Success")
		},
	)

	s.historyCard = widget.NewCard("", theme.AppSodorAddAlertGroupHistory,
		container.NewScroll(container.NewBorder(s.historyHeader, nil, nil, nil, s.historyList)))
}
