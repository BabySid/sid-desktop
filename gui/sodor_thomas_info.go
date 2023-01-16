package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/proto/sodor"
	"image/color"
	"sid-desktop/backend"
	"sid-desktop/common"
	"sid-desktop/theme"
	"strconv"
	"strings"
)

type sodorThomasInfo struct {
	thomas *sodor.ThomasInfo

	host      *widget.Entry
	port      *widget.Entry
	tags      *widget.Entry
	expandBtn *widget.Button

	tagListBinding binding.StringList
	tagList        *widget.List

	tagContent   *fyne.Container
	tagBackGroup *canvas.Rectangle

	onSubmitted func()
}

func newSodorThomasInfo(thomas *sodor.ThomasInfo) *sodorThomasInfo {
	info := sodorThomasInfo{
		thomas: thomas,
	}
	return &info
}

func (s *sodorThomasInfo) show() {
	s.host = widget.NewEntry()
	s.host.Validator = validation.NewRegexp(`\S+`, theme.AppSodorThomasInfoHost+" must not be empty")
	s.port = widget.NewEntry()
	s.port.Validator = validation.NewRegexp(`\d+`, theme.AppSodorThomasInfoPort+" must not be number")
	s.tags = widget.NewEntry()
	s.tags.OnChanged = func(str string) {
		arr := strings.Split(str, common.FavorTagSep)
		_ = s.tagListBinding.Set(arr)

		if s.tagContent.Visible() {
			s.setTagContentSize()
		}
	}

	s.expandBtn = widget.NewButtonWithIcon(theme.AppSodorThomasAddThomasExpand, theme.ResourceExpandDownIcon, nil)
	s.expandBtn.OnTapped = func() {
		if s.tagContent.Visible() {
			s.tagContent.Hide()

			s.expandBtn.SetText(theme.AppSodorThomasAddThomasExpand)
			s.expandBtn.SetIcon(theme.ResourceExpandDownIcon)
		} else {
			s.tagContent.Show()
			s.setTagContentSize()
			s.expandBtn.SetText(theme.AppSodorThomasAddThomasShrink)
			s.expandBtn.SetIcon(theme.ResourceExpandUpIcon)
		}
	}

	s.tagListBinding = binding.NewStringList()
	s.tagList = widget.NewListWithData(
		s.tagListBinding,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(item binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(item.(binding.String))
		},
	)

	s.tagBackGroup = canvas.NewRectangle(color.Transparent)
	s.tagContent = container.NewMax(s.tagList, s.tagBackGroup)
	s.tagContent.Hide()

	expandContainer := container.NewBorder(s.expandBtn, nil, nil, nil, s.tagContent)

	// set data finally
	title := theme.AppSodorAddThomas
	if s.thomas != nil { // edit or remove
		s.host.SetText(s.thomas.Host)
		s.port.SetText(fmt.Sprintf("%d", s.thomas.Port))
		s.tags.SetText(strings.Join(s.thomas.Tags, common.ThomasTagSep))
		title = theme.AppSodorEditThomas
	}

	win := dialog.NewForm(
		title, theme.ConfirmText, theme.DismissText,
		[]*widget.FormItem{
			widget.NewFormItem(theme.AppSodorThomasInfoHost, s.host),
			widget.NewFormItem(theme.AppSodorThomasInfoPort, s.port),
			widget.NewFormItem(theme.AppSodorThomasInfoTags, s.tags),
			widget.NewFormItem("", expandContainer),
		},
		func(b bool) {
			if !b {
				return
			}
			t, _ := s.tagListBinding.Get()

			id := 0
			if s.thomas != nil {
				id = int(s.thomas.Id)
			}

			portValue, err := strconv.ParseInt(s.port.Text, 10, 32)
			if err != nil {
				printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
				return
			}

			req := sodor.ThomasInfo{
				Id:   int32(id),
				Host: s.host.Text,
				Port: int32(portValue),
				Tags: t,
			}

			resp := sodor.ThomasReply{}

			if s.thomas != nil {
				if err = backend.GetSodorClient().Call(backend.UpdateThomas, &req, &resp); err != nil {
					printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
				}
			} else {
				if err = backend.GetSodorClient().Call(backend.AddThomas, &req, &resp); err != nil {
					printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
				}
			}

			if err == nil && s.onSubmitted != nil {
				s.onSubmitted()
			}
		},
		globalWin.win,
	)

	win.Resize(fyne.NewSize(400, 200))

	win.Show()
}

func (s *sodorThomasInfo) setTagContentSize() {
	size := s.tagListBinding.Length()
	if size >= 3 {
		size = 3
	}

	s.tagBackGroup.SetMinSize(fyne.NewSize(s.expandBtn.Size().Width, common.GetItemsHeightInList(s.tagList, size)))
}
