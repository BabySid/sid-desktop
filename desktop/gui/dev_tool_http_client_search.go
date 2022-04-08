package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"sid-desktop/desktop/common"
	"sid-desktop/desktop/storage"
	sidTheme "sid-desktop/desktop/theme"
	sidWidget "sid-desktop/desktop/widget"
	"strings"
)

type devToolHttpClientSearch struct {
	httpClient *devToolHttpClient

	searchText *sidWidget.CompletionEntry

	// preview
	preMethod        *widget.Label
	preUrl           *widget.Label
	preHeaderBinding binding.UntypedList
	preHeader        *widget.List
	preBody          *widget.Entry

	ok      *widget.Button
	dismiss *widget.Button

	history    *common.HttpRequestList
	curRequest *common.HttpRequest

	win fyne.Window
}

func newDevToolHttpClientSearch(client *devToolHttpClient) *devToolHttpClientSearch {
	var search devToolHttpClientSearch

	search.httpClient = client

	search.ok = widget.NewButtonWithIcon(sidTheme.ConfirmText, theme.ConfirmIcon(), search.confirm)
	search.dismiss = widget.NewButtonWithIcon(sidTheme.DismissText, theme.CancelIcon(), search.close)

	search.preMethod = widget.NewLabel(sidTheme.AppDevToolsHttpCliMethodPlaceHolder)
	search.preUrl = widget.NewLabel(sidTheme.AppDevToolsHttpCliUrlPlaceHolder)
	search.preHeaderBinding = binding.NewUntypedList()
	search.preHeader = widget.NewListWithData(
		search.preHeaderBinding,
		func() fyne.CanvasObject {
			return container.NewGridWithColumns(2, widget.NewLabel(""), widget.NewLabel(""))
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			o, _ := item.(binding.Untyped).Get()
			header := o.(*common.HttpHeader)

			key := obj.(*fyne.Container).Objects[0].(*widget.Label)
			key.SetText(header.Key)

			value := obj.(*fyne.Container).Objects[1].(*widget.Label)
			value.SetText(header.Value)
		})
	search.preBody = widget.NewMultiLineEntry()
	search.preBody.Disable()

	preview := widget.NewCard("", sidTheme.AppDevToolsCliPreviewTitle,
		container.NewBorder(
			container.NewHBox(search.preMethod, search.preUrl),
			nil, nil, nil,
			container.NewGridWithColumns(2,
				widget.NewCard("", sidTheme.AppDevToolsHttpCliHeaderTabName, search.preHeader),
				widget.NewCard("", sidTheme.AppDevToolsHttpCliBodyTabName, search.preBody))))

	search.searchText = sidWidget.NewCompletionEntry([]string{})
	search.searchText.SetPlaceHolder(sidTheme.AppDevToolsCliPreviewLoading)
	search.searchText.Disable()
	search.searchText.OnChanged = search.search
	search.searchText.OnMenuNavigation = search.preview
	search.searchText.OnSelected = search.selectRequest

	main := container.NewHSplit(container.NewBorder(search.searchText, nil, nil, nil, layout.NewSpacer()), preview)
	main.SetOffset(0.5)

	search.win = fyne.CurrentApp().NewWindow(sidTheme.AppDevToolsCliSearchTitle)
	search.win.SetContent(container.NewBorder(nil, container.NewHBox(layout.NewSpacer(), search.dismiss, search.ok), nil, nil, main))
	search.win.Resize(fyne.NewSize(800, 500))
	search.win.CenterOnScreen()

	go func() {
		var err error
		search.history, err = storage.GetAppDevToolDB().LoadHttpClientHistory()
		if err != nil {
			printErr(fmt.Errorf(sidTheme.AppDevToolsFailedFormat, err))
			return
		}
		search.searchText.SetPlaceHolder(sidTheme.AppDevToolsCliPreviewSearchPlaceHolder)
		search.searchText.Enable()
	}()

	return &search
}

func (d *devToolHttpClientSearch) close() {
	d.curRequest = nil
	d.win.Close()
}

func (d *devToolHttpClientSearch) confirm() {
	if d.curRequest != nil {
		d.httpClient.loadHttpRequest(d.curRequest)
	}
	d.curRequest = nil
	d.win.Close()
}

func (d *devToolHttpClientSearch) search(s string) {
	list := d.history.Find(s)
	if list.Len() == 0 {
		return
	}

	opt := make([]string, list.Len())
	for i := 0; i < list.Len(); i++ {
		opt[i] = list.String(i)
	}

	d.searchText.SetOptions(opt)
	d.searchText.ShowCompletion()
}

func (d *devToolHttpClientSearch) preview(s string) {
	methodAndUrl := strings.Split(s, " ")
	if len(methodAndUrl) < 2 {
		return
	}
	req, exist := d.history.Get(methodAndUrl[0], strings.Join(methodAndUrl[1:], ""))
	if !exist {
		return
	}
	d.preMethod.SetText(req.Method)
	d.preUrl.SetText(req.Url)
	d.preBody.SetText(string(req.ReqBody))

	rs := make([]interface{}, 0)
	for _, header := range req.ReqHeader {
		header := &common.HttpHeader{
			Key:   header.Key,
			Value: header.Value,
		}
		rs = append(rs, header)
	}
	d.preHeaderBinding.Set(rs)

	d.curRequest = &req
}

func (d *devToolHttpClientSearch) selectRequest(s string) {
	methodAndUrl := strings.Split(s, " ")
	if len(methodAndUrl) < 2 {
		return
	}
	req, exist := d.history.Get(methodAndUrl[0], strings.Join(methodAndUrl[1:], ""))
	if !exist {
		return
	}
	d.curRequest = &req
}
