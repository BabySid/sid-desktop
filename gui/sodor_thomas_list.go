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
	fyneTheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/proto/sodor"
	"image/color"
	"sid-desktop/backend"
	"sid-desktop/common"
	"sid-desktop/theme"
	"strconv"
	"strings"
)

type sodorThomasList struct {
	tabItem *container.TabItem

	searchEntry *widget.Entry
	newThomas   *widget.Button

	thomasHeader      *widget.List
	thomasContentList *widget.List

	thomasListBinding binding.UntypedList
	thomasListCache   *common.ThomasInfosWrapper

	viewInstanceHandle func(int32)
}

func newSodorThomasList() *sodorThomasList {
	s := sodorThomasList{}

	s.searchEntry = widget.NewEntry()
	s.searchEntry.SetPlaceHolder(theme.AppSodorThomasSearchText)
	s.searchEntry.OnChanged = s.searchThomas

	s.newThomas = widget.NewButtonWithIcon(theme.AppSodorAddThomas, theme.ResourceAddIcon, func() {
		s.addThomas()
	})

	s.thomasListBinding = binding.NewUntypedList()
	s.createThomasList()

	s.tabItem = container.NewTabItemWithIcon(theme.AppSodorThomasListName, theme.ResourceTrainIcon, nil)
	s.tabItem.Content = container.NewBorder(
		container.NewGridWithColumns(2, s.searchEntry, container.NewHBox(layout.NewSpacer(), s.newThomas)),
		nil, nil, nil,
		container.NewHScroll(container.NewBorder(s.thomasHeader, nil, nil, nil, s.thomasContentList)))

	go s.loadThomasList()
	return &s
}

func (s *sodorThomasList) GetText() string {
	return s.tabItem.Text
}

func (s *sodorThomasList) GetTabItem() *container.TabItem {
	return s.tabItem
}

func (s *sodorThomasList) searchThomas(name string) {
	if name == "" {
		if s.thomasListCache != nil {
			_ = s.thomasListBinding.Set(s.thomasListCache.AsInterfaceArray())
		}
	} else {
		if s.thomasListCache != nil {
			rs := s.thomasListCache.Find(name)
			_ = s.thomasListBinding.Set(rs.AsInterfaceArray())
		}
	}
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
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(theme.AppSodorThomasInfoID)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(theme.AppSodorThomasInfoVersion)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(theme.AppSodorThomasInfoHost)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText(theme.AppSodorThomasInfoPort)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*widget.Label).SetText(theme.AppSodorThomasInfoTags)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*widget.Label).SetText(theme.AppSodorThomasInfoStatus)
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
						widget.NewButtonWithIcon(theme.AppSodorThomasListOp2, theme.ResourceEditIcon, nil),
						widget.NewButtonWithIcon(theme.AppSodorThomasListOp3, theme.ResourceRmIcon, nil),
					)),
			)
		},
		func(data binding.DataItem, item fyne.CanvasObject) {
			o, _ := data.(binding.Untyped).Get()
			info := o.(*sodor.ThomasInfo)

			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%d", info.Id))
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(info.Version)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(info.Host)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText(fmt.Sprintf("%d", info.Port))
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*widget.Label).SetText(strings.Join(info.Tags, ","))
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*widget.Label).SetText(info.Status)

			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[5].(*fyne.Container).Objects[1].(*widget.Button).SetText(theme.AppSodorThomasListOp1)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[5].(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
				if s.viewInstanceHandle != nil {
					s.viewInstanceHandle(info.Id)
				}
			}
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[5].(*fyne.Container).Objects[2].(*widget.Button).SetText(theme.AppSodorThomasListOp2)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[5].(*fyne.Container).Objects[2].(*widget.Button).OnTapped = func() {
				s.editThomas(info)
			}
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[5].(*fyne.Container).Objects[3].(*widget.Button).SetText(theme.AppSodorThomasListOp3)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[5].(*fyne.Container).Objects[3].(*widget.Button).OnTapped = func() {
				req := sodor.ThomasInfo{
					Id: info.Id,
				}
				resp := sodor.ThomasReply{}
				if err := backend.GetSodorClient().Call(backend.DropThomas, &req, &resp); err != nil {
					printErr(err)
				}
			}
		},
	)
}

func (s *sodorThomasList) editThomas(thomas *sodor.ThomasInfo) {
	s.showThomasDialog(thomas)
}

func (s *sodorThomasList) addThomas() {
	s.showThomasDialog(nil)
}

func (s *sodorThomasList) showThomasDialog(thomas *sodor.ThomasInfo) {
	host := widget.NewEntry()
	host.Validator = validation.NewRegexp(`\S+`, theme.AppSodorThomasInfoHost+" must not be empty")
	port := widget.NewEntry()
	port.Validator = validation.NewRegexp(`\d+`, theme.AppSodorThomasInfoPort+" must not be number")
	tags := widget.NewEntry()
	expand := widget.NewButtonWithIcon(theme.AppSodorThomasAddThomasExpand, theme.ResourceExpandDownIcon, nil)

	tagArray := binding.NewStringList()
	tagList := widget.NewListWithData(
		tagArray,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(item binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(item.(binding.String))
		},
	)
	tags.OnChanged = func(s string) {
		arr := strings.Split(s, common.FavorTagSep)
		_ = tagArray.Set(arr)
	}

	title := theme.AppSodorAddThomas
	if thomas != nil { // edit or remove
		host.SetText(thomas.Host)
		port.SetText(fmt.Sprintf("%d", thomas.Port))
		tags.SetText(strings.Join(thomas.Tags, common.ThomasTagSep))
		title = theme.AppSodorEditThomas
	}

	rect := canvas.NewRectangle(color.Transparent)
	tagContent := container.NewMax(tagList, rect)
	tagContent.Hide()

	expandContainer := container.NewBorder(expand, nil, nil, nil, tagContent)

	win := dialog.NewForm(
		title, theme.ConfirmText, theme.DismissText,
		[]*widget.FormItem{
			widget.NewFormItem(theme.AppSodorThomasInfoHost, host),
			widget.NewFormItem(theme.AppSodorThomasInfoPort, port),
			widget.NewFormItem(theme.AppSodorThomasInfoTags, tags),
			widget.NewFormItem("", expandContainer),
		},
		func(b bool) {
			if b {
				t, _ := tagArray.Get()

				id := 0
				if thomas != nil {
					id = int(thomas.Id)
				}

				portValue, err := strconv.ParseInt(port.Text, 10, 32)
				if err != nil {
					printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
				}
				req := sodor.ThomasInfo{
					Id:   int32(id),
					Host: host.Text,
					Port: int32(portValue),
					Tags: t,
				}

				resp := sodor.ThomasReply{}

				if thomas != nil {
					if err = backend.GetSodorClient().Call(backend.UpdateThomas, &req, &resp); err != nil {
						printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
					}
				} else {
					if err = backend.GetSodorClient().Call(backend.AddThomas, &req, &resp); err != nil {
						printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
					}
				}

				if err == nil {
					s.loadThomasList()
				}
			}
		},
		globalWin.win,
	)

	expand.OnTapped = func() {
		if tagContent.Visible() {
			tagContent.Hide()

			expand.SetText(theme.AppSodorThomasAddThomasExpand)
			expand.SetIcon(theme.ResourceExpandDownIcon)
		} else {
			tagContent.Show()
			size := tagArray.Length()
			rect.SetMinSize(fyne.NewSize(400, getListItemHeight(tagList, size)))
			expand.SetText(theme.AppSodorThomasAddThomasShrink)
			expand.SetIcon(theme.ResourceExpandUpIcon)
		}
	}

	win.Resize(fyne.NewSize(400, 200))

	win.Show()
}

func (s *sodorThomasList) loadThomasList() {
	resp := sodor.ThomasInfos{}
	err := backend.GetSodorClient().Call(backend.ListThomas, nil, &resp)
	if err != nil {
		printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
	}

	s.thomasListCache = common.NewThomasInfosWrapper(&resp)
	s.thomasListBinding.Set(s.thomasListCache.AsInterfaceArray())
}

func getListItemHeight(list *widget.List, n int) float32 {
	height := list.CreateItem().MinSize().Height
	listHeight := float32(n)*(height+2*fyneTheme.Padding()+fyneTheme.SeparatorThicknessSize()) + 2*fyneTheme.Padding()
	return listHeight
}
