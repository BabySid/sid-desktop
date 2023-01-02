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

type sodorAlertGroupList struct {
	tabItem *container.TabItem

	newAlertGroup *widget.Button

	groupHeader *widget.List
	groupList   *widget.List

	alertGroupListBinding binding.UntypedList

	viewHistoryHandle func()
}

func newSodorAlertGroupList() *sodorAlertGroupList {
	s := sodorAlertGroupList{}

	s.newAlertGroup = widget.NewButton(theme.AppSodorAddAlertGroup, func() {
		s.addAlertGroupDialog()
	})

	s.alertGroupListBinding = binding.NewUntypedList()
	s.alertGroupListBinding.Set([]interface{}{1, 2, 3})
	s.createAlertGroupList()

	s.tabItem = container.NewTabItemWithIcon(theme.AppSodorAlertGroupTabName, theme.ResourceAlertIcon, nil)
	s.tabItem.Content = container.NewBorder(
		container.NewHBox(layout.NewSpacer(), s.newAlertGroup),
		nil, nil, nil,
		container.NewHScroll(container.NewBorder(s.groupHeader, nil, nil, nil, s.groupList)))
	return &s
}

func (s *sodorAlertGroupList) GetText() string {
	return s.tabItem.Text
}

func (s *sodorAlertGroupList) GetTabItem() *container.TabItem {
	return s.tabItem
}

func (s *sodorAlertGroupList) createAlertGroupList() {
	s.groupHeader = widget.NewList(
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
					widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{}),
				),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(theme.AppSodorCreateAlertGroupID)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(theme.AppSodorCreateAlertGroupName)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(theme.AppSodorCreateAlertGroupPlugins)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText(theme.AppSodorCreateAlertGroupCreateTime)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*widget.Label).SetText(theme.AppSodorCreateAlertGroupUpdateTime)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*widget.Label).SetText("")
		},
	)

	s.groupList = widget.NewListWithData(
		s.alertGroupListBinding,
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
						widget.NewButton(theme.AppSodorCreateAlertGroupOp1, nil),
						widget.NewButton(theme.AppSodorCreateAlertGroupOp2, nil),
						widget.NewButton(theme.AppSodorCreateAlertGroupOp3, nil),
					)),
			)
		},
		func(data binding.DataItem, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText("1")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText("group1")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText("DINGDING; Wexin")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText("2022-12-12 23:12:12")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*widget.Label).SetText("2022-12-12 23:12:12")

			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*fyne.Container).Objects[1].(*widget.Button).SetText(theme.AppSodorCreateAlertGroupOp1)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
				s.editAlertGroupDialog()
			}
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*fyne.Container).Objects[2].(*widget.Button).SetText(theme.AppSodorCreateAlertGroupOp2)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*fyne.Container).Objects[2].(*widget.Button).OnTapped = func() {
				if s.viewHistoryHandle != nil {
					s.viewHistoryHandle()
				}
			}
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*fyne.Container).Objects[3].(*widget.Button).SetText(theme.AppSodorCreateAlertGroupOp3)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*fyne.Container).Objects[3].(*widget.Button).OnTapped = func() {
				// rm
			}
		},
	)
}

func (s *sodorAlertGroupList) addAlertGroupDialog() {
	id := widget.NewLabel("id")
	name := widget.NewEntry()

	plugin := make([]*widget.CheckGroup, 0)
	for i := 0; i < 2; i++ {
		cg := widget.NewCheckGroup([]string{
			"Plugin1",
			"Plugin2",
			"Plugin3",
		}, nil)
		cg.Horizontal = true
		plugin = append(plugin, cg)
	}

	plugins := container.NewVBox()
	plugins.Add(plugin[0])
	plugins.Add(plugin[1])

	form := widget.NewForm(
		widget.NewFormItem(theme.AppSodorCreateAlertGroupID, id),
		widget.NewFormItem(theme.AppSodorCreateAlertGroupName, name),
		widget.NewFormItem(theme.AppSodorCreateAlertGroupPlugins, plugins),
	)

	cont := container.NewVBox(form)
	diag := dialog.NewCustomConfirm(theme.AppSodorAlertPluginTabName, theme.ConfirmText, theme.DismissText, cont, func(b bool) {
		if b {

		}
	}, globalWin.win)

	//diag.Resize(fyne.NewSize(500, 300))

	diag.Show()
}

func (s *sodorAlertGroupList) editAlertGroupDialog() {

}
