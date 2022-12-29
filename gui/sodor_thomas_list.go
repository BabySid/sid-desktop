package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"sid-desktop/theme"
)

type sodorThomasList struct {
	tabItem *container.TabItem

	searchEntry *widget.Entry
	newThomas   *widget.Button

	thomasHeader      *widget.List
	thomasContentList *widget.List

	thomasListBinding binding.UntypedList
}

func newSodorThomasList() *sodorThomasList {
	s := sodorThomasList{}

	s.searchEntry = widget.NewEntry()
	s.searchEntry.SetPlaceHolder(theme.AppSodorThomasSearchText)
	s.searchEntry.OnChanged = s.searchThomas

	s.newThomas = widget.NewButtonWithIcon(theme.AppSodorAddThomas, theme.ResourceAddIcon, func() {
	})

	s.thomasListBinding = binding.NewUntypedList()
	t := make([]interface{}, 10)
	for i := 0; i < 10; i++ {
		t[i] = i
	}
	s.thomasListBinding.Set(t)
	s.createThomasList()

	s.tabItem = container.NewTabItemWithIcon(theme.AppSodorThomasListName, theme.ResourceTrainIcon, nil)
	s.tabItem.Content = container.NewBorder(
		container.NewGridWithColumns(2, s.searchEntry, container.NewHBox(layout.NewSpacer(), s.newThomas)),
		nil, nil, nil,
		container.NewHScroll(container.NewBorder(s.thomasHeader, nil, nil, nil, s.thomasContentList)))
	return &s
}

func (s *sodorThomasList) GetText() string {
	return s.tabItem.Text
}

func (s *sodorThomasList) GetTabItem() *container.TabItem {
	return s.tabItem
}

func (s *sodorThomasList) searchThomas(name string) {

}

func (s *sodorThomasList) createThomasList() {
	s.thomasHeader = widget.NewList(
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
					widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{}),
				),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(theme.AppSodorThomasListHeader1)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(theme.AppSodorThomasListHeader2)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(theme.AppSodorThomasListHeader3)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText(theme.AppSodorThomasListHeader4)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*widget.Label).SetText(theme.AppSodorThomasListHeader5)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*widget.Label).SetText(theme.AppSodorThomasListHeader6)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[5].(*widget.Label).SetText("")
		},
	)

	s.thomasContentList = widget.NewListWithData(
		s.thomasListBinding,
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
						widget.NewButtonWithIcon(theme.AppSodorThomasListOp1, theme.ResourceInstanceIcon, nil),
						widget.NewButtonWithIcon(theme.AppSodorThomasListOp2, theme.ResourceRmIcon, nil),
					)),
			)
		},
		func(data binding.DataItem, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText("1")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText("thomas version")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText("host")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText("port")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*widget.Label).SetText("pid")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*widget.Label).SetText("status")

			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[5].(*fyne.Container).Objects[1].(*widget.Button).SetText(theme.AppSodorThomasListOp1)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[5].(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {

			}
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[5].(*fyne.Container).Objects[2].(*widget.Button).SetText(theme.AppSodorThomasListOp2)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[5].(*fyne.Container).Objects[2].(*widget.Button).OnTapped = func() {

			}
		},
	)
}
