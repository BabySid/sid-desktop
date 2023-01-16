package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/proto/sodor"
	"image/color"
	"sid-desktop/backend"
	"sid-desktop/common"
	"sid-desktop/theme"
	"strings"
)

type sodorThomasList struct {
	tabItem *container.TabItem

	searchEntry *widget.Entry
	newThomas   *widget.Button

	thomasHeader      fyne.CanvasObject
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

func (s *sodorThomasList) createThomasContListOpButtons() *fyne.Container {
	return container.NewHBox(
		layout.NewSpacer(),
		widget.NewButtonWithIcon(theme.AppSodorThomasListOp1, theme.ResourceInstanceIcon, nil),
		widget.NewButtonWithIcon(theme.AppSodorThomasListOp2, theme.ResourceEditIcon, nil),
		widget.NewButtonWithIcon(theme.AppSodorThomasListOp3, theme.ResourceRmIcon, nil),
	)
}

func (s *sodorThomasList) createThomasList() {
	size := s.createThomasContListOpButtons().MinSize()
	spaceLabel := canvas.NewRectangle(color.Transparent)
	spaceLabel.SetMinSize(fyne.NewSize(size.Width, size.Height))

	s.thomasHeader = container.NewBorder(nil, nil,
		widget.NewLabelWithStyle(theme.AppSodorThomasInfoID, fyne.TextAlignLeading, fyne.TextStyle{}),
		spaceLabel,
		container.NewGridWithColumns(4,
			widget.NewLabelWithStyle(theme.AppSodorThomasInfoVersion, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(theme.AppSodorThomasInfoHostPort, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(theme.AppSodorThomasInfoTags, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(theme.AppSodorThomasInfoStatus, fyne.TextAlignCenter, fyne.TextStyle{}),
		),
	)

	s.thomasContentList = widget.NewListWithData(
		s.thomasListBinding,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil,
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				s.createThomasContListOpButtons(),
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
			info := o.(*sodor.ThomasInfo)

			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%d", info.Id))
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(info.Version)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%s:%d", info.Host, info.Port))
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText(strings.Join(info.Tags, ","))
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*widget.Label).SetText(info.Status)

			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[1].(*widget.Button).SetText(theme.AppSodorThomasListOp1)
			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
				if s.viewInstanceHandle != nil {
					s.viewInstanceHandle(info.Id)
				}
			}
			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[2].(*widget.Button).SetText(theme.AppSodorThomasListOp2)
			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[2].(*widget.Button).OnTapped = func() {
				s.editThomas(info)
			}
			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[3].(*widget.Button).SetText(theme.AppSodorThomasListOp3)
			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[3].(*widget.Button).OnTapped = func() {
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
	info := newSodorThomasInfo(thomas)
	info.onSubmitted = s.loadThomasList
	info.show()
}

func (s *sodorThomasList) loadThomasList() {
	resp := sodor.ThomasInfos{}
	resp.ThomasInfos = make([]*sodor.ThomasInfo, 0)
	//err := backend.GetSodorClient().Call(backend.ListThomas, nil, &resp)
	//if err != nil {
	//	printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
	//	return
	//}

	for i := 0; i < 3; i++ {
		t := sodor.ThomasInfo{}
		t.Id = int32(i)
		t.Host = "127.0.0.1"
		t.Port = 12345
		t.Status = "OK"
		t.Tags = []string{"hello", "world", "gogogo"}
		t.Version = "thomas_2022-12-23_16:00:21"
		resp.ThomasInfos = append(resp.ThomasInfos, &t)
	}
	s.thomasListCache = common.NewThomasInfosWrapper(&resp)
	s.thomasListBinding.Set(s.thomasListCache.AsInterfaceArray())
}
