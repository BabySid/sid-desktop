package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	"github.com/BabySid/proto/sodor"
	"sid-desktop/common"
	"sid-desktop/theme"
	"strings"
)

var _ sodorInterface = (*sodorAlertPlugin)(nil)

type sodorAlertPlugin struct {
	sodorAdapter

	docs *container.AppTabs

	addInstance          *widget.Button
	refresh              *widget.Button
	instanceListBinding  binding.UntypedList
	pluginInstanceHeader fyne.CanvasObject
	pluginInstance       *widget.List
}

func (s *sodorAlertPlugin) CreateView() fyne.CanvasObject {
	if s.content != nil {
		return s.content
	}

	s.addInstance = widget.NewButton(theme.AppSodorCreateAlertPluginInstance, func() {
		s.addPlugin()
	})
	s.refresh = widget.NewButton(theme.AppPageRefresh, func() {
		s.loadAlertPlugins()
	})

	s.pluginInstanceHeader = container.NewBorder(nil, nil,
		widget.NewLabelWithStyle(theme.AppSodorCreateAlertPluginInstanceID, fyne.TextAlignLeading, fyne.TextStyle{}),
		nil,
		container.NewGridWithColumns(5,
			widget.NewLabelWithStyle(theme.AppSodorCreateAlertPluginInstanceName, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(theme.AppSodorCreateAlertPluginInstancePlugin, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(theme.AppSodorCreateAlertPluginInstanceCreateTime, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(theme.AppSodorCreateAlertPluginInstanceUpdateTime, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{})),
	)

	s.instanceListBinding = binding.NewUntypedList()
	s.pluginInstance = widget.NewListWithData(
		s.instanceListBinding,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil,
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				nil,
				container.NewGridWithColumns(5,
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					container.NewHBox(
						layout.NewSpacer(),
						widget.NewButton(theme.AppSodorCreateAlertPluginOp1, nil),
						widget.NewButton(theme.AppSodorCreateAlertPluginOp2, nil),
					)),
			)
		},
		func(data binding.DataItem, item fyne.CanvasObject) {
			o, _ := data.(binding.Untyped).Get()
			plugin := o.(*sodor.AlertPluginInstance)

			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%d", plugin.Id))
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(plugin.Name)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(plugin.PluginName)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText(gobase.FormatTimeStamp(int64(plugin.CreateAt)))
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*widget.Label).SetText(gobase.FormatTimeStamp(int64(plugin.UpdateAt)))

			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*fyne.Container).Objects[1].(*widget.Button).SetText(theme.AppSodorCreateAlertPluginOp1)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
				s.editPlugin(plugin)
			}
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*fyne.Container).Objects[2].(*widget.Button).SetText(theme.AppSodorCreateAlertPluginOp2)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*fyne.Container).Objects[2].(*widget.Button).OnTapped = func() {
				req := sodor.AlertPluginInstance{
					Id: plugin.Id,
				}
				resp := sodor.AlertPluginReply{}
				if err := common.GetSodorClient().Call(common.DeleteAlertPluginInstance, &req, &resp); err != nil {
					printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
					return
				}
				s.loadAlertPlugins()
			}
		},
	)

	s.docs = container.NewAppTabs()
	s.docs.Append(
		container.NewTabItem(theme.AppSodorAlertPluginTabName,
			container.NewBorder(
				container.NewHBox(layout.NewSpacer(), s.refresh, s.addInstance),
				nil, nil, nil,
				container.NewBorder(s.pluginInstanceHeader, nil, nil, nil, s.pluginInstance))),
	)
	s.docs.SetTabLocation(container.TabLocationTop)

	s.content = s.docs

	go s.loadAlertPlugins()

	return s.content
}

func (s *sodorAlertPlugin) loadAlertPlugins() {
	resp := sodor.AlertPluginInstances{}
	err := common.GetSodorClient().Call(common.ListAlertPluginInstances, nil, &resp)
	if err != nil {
		printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
		return
	}

	wrapper := common.NewSodorAlertPluginsWrapper(&resp)
	s.instanceListBinding.Set(wrapper.AsInterfaceArray())
}

func (s *sodorAlertPlugin) addPlugin() {
	diag := s.buildPluginInstanceDialog(nil)
	diag.Show()
}

func (s *sodorAlertPlugin) editPlugin(plugin *sodor.AlertPluginInstance) {
	diag := s.buildPluginInstanceDialog(plugin)
	diag.Show()
}

func (s *sodorAlertPlugin) buildPluginInstanceDialog(plugin *sodor.AlertPluginInstance) dialog.Dialog {
	var diagCont []*widget.FormItem

	name := widget.NewEntry()
	name.Validator = validation.NewRegexp(`\S+`, theme.AppSodorCreateAlertPluginDingDingWebHook+" must not be empty")

	webhook := widget.NewEntry()
	webhook.Validator = validation.NewRegexp(`\S+`, theme.AppSodorCreateAlertPluginDingDingWebHook+" must not be empty")
	sign := widget.NewEntry()
	sign.Validator = validation.NewRegexp(`\S+`, theme.AppSodorCreateAlertPluginDingDingSign+" must not be empty")
	atMobiles := widget.NewEntry()

	plugins := widget.NewSelect([]string{sodor.AlertPluginName_APN_DingDing.String()}, func(s string) {
		if s == sodor.AlertPluginName_APN_DingDing.String() {
			webhook.Show()
			sign.Show()
			atMobiles.Show()
		} else {
			webhook.Hide()
			sign.Hide()
			atMobiles.Hide()
		}
	})
	plugins.SetSelectedIndex(0)

	title := theme.AppSodorCreateAlertPluginInstance
	if plugin == nil {
		diagCont = []*widget.FormItem{
			widget.NewFormItem(theme.AppSodorCreateAlertPluginInstanceName, name),
			widget.NewFormItem(theme.AppSodorCreateAlertPluginInstancePlugin, plugins),
		}
	} else {
		title = theme.AppSodorEditAlertPluginInstance

		name.SetText(plugin.Name)
		webhook.SetText(plugin.Dingding.Webhook)
		sign.SetText(plugin.Dingding.Sign)
		atMobiles.SetText(strings.Join(plugin.Dingding.AtMobiles, common.ArraySeparator))

		id := widget.NewLabel(fmt.Sprintf("%d", plugin.Id))
		diagCont = []*widget.FormItem{
			widget.NewFormItem(theme.AppSodorCreateAlertPluginInstanceID, id),
			widget.NewFormItem(theme.AppSodorCreateAlertPluginInstanceName, name),
			widget.NewFormItem(theme.AppSodorCreateAlertPluginInstancePlugin, plugins),
		}
	}

	diagCont = append(diagCont, widget.NewFormItem(theme.AppSodorCreateAlertPluginDingDingWebHook, webhook))
	diagCont = append(diagCont, widget.NewFormItem(theme.AppSodorCreateAlertPluginDingDingSign, sign))
	diagCont = append(diagCont, widget.NewFormItem(theme.AppSodorCreateAlertPluginDingDingAtMobiles, atMobiles))

	diag := dialog.NewForm(title, theme.ConfirmText, theme.DismissText,
		diagCont, func(b bool) {
			if !b {
				return
			}

			id := 0
			if plugin != nil {
				id = int(plugin.Id)
			}

			req := sodor.AlertPluginInstance{
				Id:         int32(id),
				Name:       name.Text,
				PluginName: plugins.Selected,
				Dingding: &sodor.AlertPluginDingDing{
					Webhook:   webhook.Text,
					Sign:      sign.Text,
					AtMobiles: gobase.SplitAndTrimSpace(atMobiles.Text, common.ArraySeparator),
				},
			}
			resp := sodor.AlertPluginReply{}

			var err error
			if plugin != nil {
				if err = common.GetSodorClient().Call(common.UpdateAlertPluginInstance, &req, &resp); err != nil {
					printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
				}
			} else {
				if err = common.GetSodorClient().Call(common.CreateAlertPluginInstance, &req, &resp); err != nil {
					printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
				}
			}

			if err == nil {
				s.loadAlertPlugins()
			}
		}, globalWin.win)

	diag.Resize(fyne.NewSize(500, 300))
	return diag
}
