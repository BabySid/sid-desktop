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
	"strconv"
	"strings"
)

type sodorAlertGroupHistory struct {
	group *sodor.AlertGroup

	refresh *widget.Button

	tabItem *container.TabItem

	groupID             *widget.Label
	groupName           *widget.Label
	pluginInstanceNames *widget.Select
	headerContainer     *fyne.Container

	historyCard        *widget.Card
	historyListBinding binding.UntypedList
	historyHeader      fyne.CanvasObject
	historyList        *widget.List
}

func newSodorAlertGroupHistory(group *sodor.AlertGroup) *sodorAlertGroupHistory {
	ins := sodorAlertGroupHistory{}
	ins.group = group

	ins.buildHistoryInfo()
	ins.buildPluginNames()

	ins.tabItem = container.NewTabItem(theme.AppSodorAddAlertGroupHistory, nil)
	ins.tabItem.Content = container.NewBorder(
		ins.headerContainer, nil, nil, nil,
		ins.historyCard)

	return &ins
}

func (s *sodorAlertGroupHistory) buildPluginNames() {
	s.groupID = widget.NewLabel(fmt.Sprintf("%d", s.group.Id))
	s.groupName = widget.NewLabel(s.group.Name)

	s.refresh = widget.NewButtonWithIcon(theme.AppPageRefresh, theme.ResourceRefreshIcon, func() {
		v := s.pluginInstanceNames.Selected
		rs := strings.SplitN(v, ":", 2)
		id, _ := strconv.Atoi(rs[0])
		s.loadAlertGroupHistory(int32(id))
	})

	plugins := common.GetSodorCache().GetAlertPluginInstances(s.group.PluginInstances...)
	opts := make([]string, 0)
	for _, p := range plugins.AlertPluginInstances {
		opts = append(opts, fmt.Sprintf("%d:%s", p.Id, p.Name))
	}

	s.pluginInstanceNames = widget.NewSelect(opts, func(v string) {
		rs := strings.SplitN(v, ":", 2)
		id, _ := strconv.Atoi(rs[0])
		s.loadAlertGroupHistory(int32(id))
	})
	s.pluginInstanceNames.SetSelectedIndex(0)

	s.headerContainer = container.NewVBox(
		container.NewHBox(
			widget.NewForm(widget.NewFormItem(theme.AppSodorCreateAlertGroupID, s.groupID)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorCreateAlertGroupName, s.groupName)),
			s.pluginInstanceNames,
			layout.NewSpacer(),
			s.refresh),
	)
}

func (s *sodorAlertGroupHistory) buildHistoryInfo() {
	s.historyListBinding = binding.NewUntypedList()

	s.historyHeader = container.NewBorder(nil, nil,
		widget.NewLabelWithStyle(theme.AppSodorCreateAlertGroupPluginID, fyne.TextAlignLeading, fyne.TextStyle{}),
		widget.NewLabelWithStyle(theme.AppSodorCreateAlertGroupStatusMsg, fyne.TextAlignCenter, fyne.TextStyle{}),
		container.NewBorder(nil, nil, nil,
			widget.NewLabelWithStyle(theme.AppSodorCreateAlertGroupAlertCreateTime, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(theme.AppSodorCreateAlertGroupAlertMsg, fyne.TextAlignCenter, fyne.TextStyle{}),
		),
	)

	s.historyList = widget.NewListWithData(
		s.historyListBinding,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil,
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
				container.NewBorder(nil, nil, nil,
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
				),
			)
		},
		func(data binding.DataItem, item fyne.CanvasObject) {
			o, _ := data.(binding.Untyped).Get()
			his := o.(*sodor.AlertPluginInstanceHistory)

			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%d", his.Id))
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(his.AlertMsg)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(gobase.FormatTimeStamp(int64(his.CreateAt)))
			item.(*fyne.Container).Objects[2].(*widget.Label).SetText(his.StatusMsg)
		},
	)

	s.historyCard = widget.NewCard("", theme.AppSodorAddAlertGroupHistory,
		container.NewScroll(container.NewBorder(s.historyHeader, nil, nil, nil, s.historyList)))
}

func (s *sodorAlertGroupHistory) loadAlertGroupHistory(pluginID int32) {
	resp := sodor.AlertPluginInstanceHistories{}
	req := sodor.AlertPluginInstanceHistory{}
	req.GroupId = s.group.Id
	req.InstanceId = pluginID
	err := common.GetSodorClient().Call(common.ShowAlertPluginInstanceHistories, &req, &resp)
	if err != nil {
		printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
		return
	}

	data := common.NewSodorAlertGroupHistoriesWrapperWrapper(&resp)
	s.historyListBinding.Set(data.AsInterfaceArray())
}
