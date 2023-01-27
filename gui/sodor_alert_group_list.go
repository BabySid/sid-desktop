package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/proto/sodor"
	"image/color"
	"sid-desktop/common"
	"sid-desktop/theme"
	"strconv"
	"strings"
)

type sodorAlertGroupList struct {
	tabItem *container.TabItem

	refresh       *widget.Button
	newAlertGroup *widget.Button

	groupHeader fyne.CanvasObject
	groupList   *widget.List

	alertGroupListBinding binding.UntypedList

	viewHistoryHandle func(group *sodor.AlertGroup)
}

func newSodorAlertGroupList() *sodorAlertGroupList {
	s := sodorAlertGroupList{}

	s.refresh = widget.NewButton(theme.AppPageRefresh, func() {
		s.loadAlertGroupList()
	})
	s.newAlertGroup = widget.NewButton(theme.AppSodorAddAlertGroup, func() {
		s.addAlertGroupDialog()
	})

	s.createAlertGroupList()

	s.tabItem = container.NewTabItemWithIcon(theme.AppSodorAlertGroupTabName, theme.ResourceAlertIcon, nil)
	s.tabItem.Content = container.NewBorder(
		container.NewHBox(layout.NewSpacer(), s.refresh, s.newAlertGroup),
		nil, nil, nil,
		container.NewBorder(s.groupHeader, nil, nil, nil, s.groupList))

	go s.loadAlertGroupList()
	return &s
}

func (s *sodorAlertGroupList) GetText() string {
	return s.tabItem.Text
}

func (s *sodorAlertGroupList) GetTabItem() *container.TabItem {
	return s.tabItem
}

func (s *sodorAlertGroupList) createAlertGroupOpButtons() *fyne.Container {
	return container.NewHBox(
		widget.NewButton(theme.AppSodorCreateAlertGroupOp1, nil),
		widget.NewButton(theme.AppSodorCreateAlertGroupOp2, nil),
		widget.NewButton(theme.AppSodorCreateAlertGroupOp3, nil),
	)
}

func (s *sodorAlertGroupList) createAlertGroupList() {
	size := s.createAlertGroupOpButtons().MinSize()
	spaceLabel := canvas.NewRectangle(color.Transparent)
	spaceLabel.SetMinSize(fyne.NewSize(size.Width, size.Height))

	s.groupHeader = container.NewBorder(nil, nil,
		widget.NewLabelWithStyle(theme.AppSodorCreateAlertGroupID, fyne.TextAlignLeading, fyne.TextStyle{}),
		spaceLabel,
		container.NewGridWithColumns(4,
			widget.NewLabelWithStyle(theme.AppSodorCreateAlertGroupName, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(theme.AppSodorCreateAlertGroupPlugins, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(theme.AppSodorCreateAlertGroupCreateTime, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(theme.AppSodorCreateAlertGroupUpdateTime, fyne.TextAlignCenter, fyne.TextStyle{}),
		),
	)

	s.alertGroupListBinding = binding.NewUntypedList()
	s.groupList = widget.NewListWithData(
		s.alertGroupListBinding,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil,
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				s.createAlertGroupOpButtons(),
				container.NewGridWithColumns(4,
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
				),
			)
		},
		func(data binding.DataItem, item fyne.CanvasObject) {
			o, _ := data.(binding.Untyped).Get()
			group := o.(common.SodorAlertGroup)

			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(group.ID)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(group.Name)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(group.PluginNames)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText(group.CreateTime)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*widget.Label).SetText(group.UpdateTime)

			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[0].(*widget.Button).SetText(theme.AppSodorCreateAlertGroupOp1)
			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[0].(*widget.Button).OnTapped = func() {
				s.editAlertGroupDialog(group.GroupObj)
			}
			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[1].(*widget.Button).SetText(theme.AppSodorCreateAlertGroupOp2)
			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
				if s.viewHistoryHandle != nil {
					s.viewHistoryHandle(group.GroupObj)
				}
			}
			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[2].(*widget.Button).SetText(theme.AppSodorCreateAlertGroupOp3)
			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[2].(*widget.Button).OnTapped = func() {
				req := sodor.AlertGroup{
					Id: group.GroupObj.Id,
				}
				resp := sodor.AlertGroupReply{}
				if err := common.GetSodorClient().Call(common.DeleteAlertGroup, &req, &resp); err != nil {
					printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
					return
				}
				s.loadAlertGroupList()
			}
		},
	)
}

func (s *sodorAlertGroupList) addAlertGroupDialog() {
	if !s.checkBeforeBuild() {
		return
	}
	diag := s.buildAlertGroupDialog(nil)
	diag.Show()
}

func (s *sodorAlertGroupList) editAlertGroupDialog(group *sodor.AlertGroup) {
	if !s.checkBeforeBuild() {
		return
	}
	diag := s.buildAlertGroupDialog(group)
	diag.Show()
}

func (s *sodorAlertGroupList) checkBeforeBuild() bool {
	pluginIns := common.GetSodorCache().GetAlertPluginInstances()
	if pluginIns == nil || len(pluginIns.AlertPluginInstances) == 0 {
		printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, theme.AppSodorEmptyAlertPluginInstance))
		return false
	}

	return true
}

func (s *sodorAlertGroupList) buildAlertGroupDialog(group *sodor.AlertGroup) dialog.Dialog {
	var diagCont []*widget.FormItem

	name := widget.NewEntry()
	name.Validator = validation.NewRegexp(`\S+`, theme.AppSodorCreateAlertGroupName+" must not be empty")
	nameFormItem := &widget.FormItem{
		Text:     theme.AppSodorCreateAlertGroupName,
		Widget:   name,
		HintText: theme.AppSodorAlertGroupNameTooltip,
	}

	// init plugin instance
	pluginIns := common.GetSodorCache().GetAlertPluginInstances()

	opts := make([]string, len(pluginIns.AlertPluginInstances))
	for i := 0; i < len(pluginIns.AlertPluginInstances); i++ {
		p := pluginIns.AlertPluginInstances[len(pluginIns.AlertPluginInstances)-1-i]
		opts[i] = fmt.Sprintf("%d:%s", p.Id, p.Name)
	}

	plugin := widget.NewCheckGroup(opts, nil)
	plugin.Required = true

	title := theme.AppSodorAddAlertGroup
	selOpts := make([]string, 0)
	if group == nil {
		diagCont = []*widget.FormItem{nameFormItem}
		if len(opts) > 0 {
			selOpts = append(selOpts, opts[0])
		}
	} else {
		title = theme.AppSodorEditAlertGroup
		name.SetText(group.Name)
		id := widget.NewLabel(fmt.Sprintf("%d", group.Id))

		for _, insID := range group.PluginInstances {
			for _, o := range opts {
				rs := strings.SplitN(o, ":", 2)
				if fmt.Sprintf("%d", insID) == rs[0] {
					selOpts = append(selOpts, o)
				}
			}
		}

		diagCont = []*widget.FormItem{
			widget.NewFormItem(theme.AppSodorCreateAlertGroupID, id),
			nameFormItem,
		}
	}
	plugin.SetSelected(selOpts)
	back := canvas.NewRectangle(color.Transparent)
	back.SetMinSize(s.setPluginOptionContentSize(opts))
	diagCont = append(diagCont, widget.NewFormItem(theme.AppSodorCreateAlertGroupPlugins, container.NewMax(
		container.NewMax(container.NewScroll(plugin), back))))

	diag := dialog.NewForm(title, theme.ConfirmText, theme.DismissText,
		diagCont, func(b bool) {
			if !b {
				return
			}

			pluginIDs := make([]int32, 0)
			for _, opt := range plugin.Selected {
				rs := strings.SplitN(opt, ":", 2)
				id, _ := strconv.Atoi(rs[0])
				pluginIDs = append(pluginIDs, int32(id))
			}

			id := 0
			if group != nil {
				id = int(group.Id)
			}
			req := sodor.AlertGroup{
				Id:              int32(id),
				Name:            name.Text,
				PluginInstances: pluginIDs,
			}
			resp := sodor.AlertGroupReply{}

			var err error
			if group != nil {
				if err = common.GetSodorClient().Call(common.UpdateAlertGroup, &req, &resp); err != nil {
					printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
				}
			} else {
				if err = common.GetSodorClient().Call(common.CreateAlertGroup, &req, &resp); err != nil {
					printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
				}
			}

			if err == nil {
				s.loadAlertGroupList()
			}
		}, globalWin.win)

	diag.Resize(fyne.NewSize(500, 300))
	return diag
}

func (s *sodorAlertGroupList) loadAlertGroupList() {
	err := common.GetSodorCache().LoadAlertGroups()
	if err != nil {
		printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
		return
	}

	wrapper := common.NewSodorAlertGroupsWrapperWrapper(common.GetSodorCache().GetAlertGroups())
	s.alertGroupListBinding.Set(wrapper.AsInterfaceArray())
}

func (s *sodorAlertGroupList) setPluginOptionContentSize(opts []string) fyne.Size {
	size := len(opts)
	if size >= 3 {
		size = 3
	}

	return fyne.NewSize(200, common.GetItemsHeightInCheck(size))
}
