package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"sid-desktop/theme"
)

var _ sodorInterface = (*sodorAlertPlugin)(nil)

type sodorAlertPlugin struct {
	sodorAdapter

	docs *container.AppTabs

	addInstance          *widget.Button
	instanceListBinding  binding.UntypedList
	pluginInstanceHeader *widget.List
	pluginInstance       *widget.List
}

func (s *sodorAlertPlugin) CreateView() fyne.CanvasObject {
	if s.content != nil {
		return s.content
	}

	s.addInstance = widget.NewButton(theme.AppSodorCreateAlertPluginInstance, func() {
		s.addPlugin()
	})

	s.pluginInstanceHeader = widget.NewList(
		func() int {
			return 1
		},
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil,
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				nil,
				container.NewGridWithColumns(5,
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{})),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(theme.AppSodorCreateAlertPluginInstanceID)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(theme.AppSodorCreateAlertPluginInstanceName)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(theme.AppSodorCreateAlertPluginInstancePlugin)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText(theme.AppSodorCreateAlertPluginInstanceCreateTime)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*widget.Label).SetText(theme.AppSodorCreateAlertPluginInstanceUpdateTime)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*widget.Label).SetText("")
		},
	)

	s.instanceListBinding = binding.NewUntypedList()
	s.instanceListBinding.Set([]interface{}{1, 2, 3})
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
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText("1")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText("instance name")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText("APN_DINGDING")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText("2022-12-12 23:22:15")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*widget.Label).SetText("2022-12-15 23:22:15")

			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*fyne.Container).Objects[1].(*widget.Button).SetText(theme.AppSodorCreateAlertPluginOp1)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
				s.editPlugin()
			}
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*fyne.Container).Objects[2].(*widget.Button).SetText(theme.AppSodorCreateAlertPluginOp2)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*fyne.Container).Objects[2].(*widget.Button).OnTapped = func() {
				s.rmPlugin()
			}
		},
	)

	s.docs = container.NewAppTabs()
	s.docs.Append(
		container.NewTabItem(theme.AppSodorAlertPluginTabName,
			container.NewBorder(
				container.NewHBox(layout.NewSpacer(), s.addInstance),
				nil, nil, nil,
				container.NewScroll(container.NewBorder(s.pluginInstanceHeader, nil, nil, nil, s.pluginInstance)))),
	)
	s.docs.SetTabLocation(container.TabLocationTop)

	s.content = s.docs
	return s.content
}

func (s *sodorAlertPlugin) addPlugin() {
	diag := s.buildPluginInstanceDialog()
	diag.Show()
}

func (s *sodorAlertPlugin) editPlugin() {
	diag := s.buildPluginInstanceDialog()
	diag.Show()
}

func (s *sodorAlertPlugin) rmPlugin() {

}

func (s *sodorAlertPlugin) buildPluginInstanceDialog() dialog.Dialog {
	id := widget.NewLabel("id")
	name := widget.NewEntry()

	webhook := widget.NewEntry()
	sign := widget.NewEntry()
	atMobiles := widget.NewEntry()
	dingding := widget.NewCard("", theme.AppSodorCreateAlertPluginDingDing,
		widget.NewForm(
			widget.NewFormItem(theme.AppSodorCreateAlertPluginDingDingWebHook, webhook),
			widget.NewFormItem(theme.AppSodorCreateAlertPluginDingDingSign, sign),
			widget.NewFormItem(theme.AppSodorCreateAlertPluginDingDingAtMobiles, atMobiles),
		),
	)

	plugin := widget.NewSelect([]string{"Dingding", "Weixin"}, func(s string) {
		if s == "Dingding" {
			dingding.Show()
		} else {
			dingding.Hide()
		}
	})
	plugin.SetSelectedIndex(0)

	form := widget.NewForm(
		widget.NewFormItem(theme.AppSodorCreateAlertPluginInstanceID, id),
		widget.NewFormItem(theme.AppSodorCreateAlertPluginInstanceName, name),
		widget.NewFormItem(theme.AppSodorCreateAlertPluginInstancePlugin, plugin),
	)

	cont := container.NewVBox(form, dingding)
	diag := dialog.NewCustomConfirm(theme.AppSodorAlertPluginTabName, theme.ConfirmText, theme.DismissText, cont, func(b bool) {
		if b {

		}
	}, globalWin.win)

	diag.Resize(fyne.NewSize(500, 500))

	return diag
}
