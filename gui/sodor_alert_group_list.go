package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	"github.com/BabySid/proto/sodor"
	"image/color"
	"sid-desktop/theme"
)

type sodorAlertGroupList struct {
	tabItem *container.TabItem

	refresh       *widget.Button
	newAlertGroup *widget.Button

	groupHeader fyne.CanvasObject
	groupList   *widget.List

	alertGroupListBinding binding.UntypedList

	viewHistoryHandle func()
}

func newSodorAlertGroupList() *sodorAlertGroupList {
	s := sodorAlertGroupList{}

	s.refresh = widget.NewButton(theme.AppPageRefresh, func() {
		s.loadAlertGroupList()
	})
	s.newAlertGroup = widget.NewButton(theme.AppSodorAddAlertGroup, func() {
		s.addAlertGroupDialog()
	})

	s.alertGroupListBinding = binding.NewUntypedList()
	s.createAlertGroupList()

	s.tabItem = container.NewTabItemWithIcon(theme.AppSodorAlertGroupTabName, theme.ResourceAlertIcon, nil)
	s.tabItem.Content = container.NewBorder(
		container.NewHBox(layout.NewSpacer(), s.refresh, s.newAlertGroup),
		nil, nil, nil,
		container.NewBorder(s.groupHeader, nil, nil, nil, s.groupList))
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
			group := o.(*sodor.AlertGroup)

			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%d", group.Id))
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(group.Name)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText("")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText(gobase.FormatTimeStamp(int64(group.CreateAt)))
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*widget.Label).SetText(gobase.FormatTimeStamp(int64(group.UpdateAt)))

			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[0].(*widget.Button).SetText(theme.AppSodorCreateAlertGroupOp1)
			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[0].(*widget.Button).OnTapped = func() {
				s.editAlertGroupDialog(group)
			}
			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[1].(*widget.Button).SetText(theme.AppSodorCreateAlertGroupOp2)
			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
				if s.viewHistoryHandle != nil {
					s.viewHistoryHandle()
				}
			}
			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[2].(*widget.Button).SetText(theme.AppSodorCreateAlertGroupOp3)
			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[2].(*widget.Button).OnTapped = func() {
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

func (s *sodorAlertGroupList) editAlertGroupDialog(group *sodor.AlertGroup) {

}

func (s *sodorAlertGroupList) loadAlertGroupList() {

}
