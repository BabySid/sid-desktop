package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"sid-desktop/desktop/common"
	"strings"
)

var _ devToolInterface = (*devToolHttpClient)(nil)

type devToolHttpClient struct {
	method       *widget.Select
	url          *widget.Entry
	sendRequest  *widget.Button
	hideRequest  *widget.Button
	hideResponse *widget.Button

	reqBodyArea      *container.AppTabs
	reqHeaderBinding binding.UntypedList
	requestHeader    *widget.List
	requestBody      *widget.Entry
	requestBodyType  *widget.RadioGroup

	respBodyArea      *container.AppTabs
	respHeaderBinding binding.UntypedList
	responseHeader    *widget.List
	responseBody      *widget.Entry
	responseBodyType  *widget.RadioGroup

	content fyne.CanvasObject
}

func (d *devToolHttpClient) CreateView() fyne.CanvasObject {
	if d.content != nil {
		return d.content
	}

	d.method = widget.NewSelect(common.HttpMethod, nil)
	d.method.PlaceHolder = d.method.Options[0]
	d.method.SetSelectedIndex(0)
	d.url = widget.NewEntry()
	d.url.SetPlaceHolder("url")
	d.sendRequest = widget.NewButton("Send", nil)
	d.hideRequest = widget.NewButton("Hide Request", func() {
		if d.reqBodyArea.Visible() {
			d.reqBodyArea.Hide()
			d.hideRequest.SetText("Show Request")
		} else {
			d.reqBodyArea.Show()
			d.hideRequest.SetText("Hide Request")
		}
	})
	d.hideResponse = widget.NewButton("Hide Response", func() {
		if d.respBodyArea.Visible() {
			d.respBodyArea.Hide()
			d.hideResponse.SetText("Show Request")
		} else {
			d.respBodyArea.Show()
			d.hideResponse.SetText("Hide Request")
		}
	})

	d.createRequestView()
	d.createResponseView()

	area := container.NewVSplit(d.reqBodyArea, d.respBodyArea)
	area.SetOffset(0.1)

	d.content = container.NewBorder(
		container.NewBorder(nil, nil, d.method, container.NewHBox(d.sendRequest, d.hideRequest, d.hideResponse), d.url),
		nil, nil, nil,
		area)

	return d.content
}

func (d *devToolHttpClient) createRequestView() {
	var req common.HttpRequest
	common.InitHttpRequest(&req)

	d.reqHeaderBinding = binding.NewUntypedList()
	d.reqHeaderBinding.Set(req.AsInterfaceArray())

	d.requestHeader = widget.NewListWithData(
		d.reqHeaderBinding,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil, widget.NewButton("Remove", nil),
				container.NewGridWithColumns(2,
					widget.NewSelectEntry(common.HttpHeaderName),
					widget.NewEntry(),
				))
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			o, _ := item.(binding.Untyped).Get()
			header := o.(common.HttpHeader)

			key := obj.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.SelectEntry)
			key.SetText(header.Key)

			value := obj.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Entry)
			value.SetText(fmt.Sprintf("%v", header.Value))

			rm := obj.(*fyne.Container).Objects[1].(*widget.Button)
			rm.OnTapped = func() {
				fmt.Println(header)
			}
		},
	)

	d.requestBodyType = widget.NewRadioGroup([]string{
		"none",
		"json",
	}, nil)
	d.requestBodyType.Horizontal = true

	d.requestBody = widget.NewMultiLineEntry()

	d.reqBodyArea = container.NewAppTabs(
		container.NewTabItem("Header", d.requestHeader),
		container.NewTabItem("Body",
			container.NewBorder(d.requestBodyType, nil, nil, nil, d.requestBody)),
	)
}

func (d *devToolHttpClient) createResponseView() {
	var req common.HttpRequest
	common.InitHttpRequest(&req)

	d.respHeaderBinding = binding.NewUntypedList()
	d.respHeaderBinding.Set(req.AsInterfaceArray())

	d.responseHeader = widget.NewListWithData(
		d.reqHeaderBinding,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil, nil,
				container.NewGridWithColumns(2,
					widget.NewSelectEntry(common.HttpHeaderName),
					widget.NewEntry(),
				))
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			o, _ := item.(binding.Untyped).Get()
			header := o.(common.HttpHeader)

			key := obj.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.SelectEntry)
			key.SetText(header.Key)

			value := obj.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Entry)
			value.SetText(fmt.Sprintf("%v", header.Value))
		},
	)

	d.responseBodyType = widget.NewRadioGroup([]string{
		"none",
		"json",
	}, nil)

	d.responseBodyType.Horizontal = true

	d.responseBody = widget.NewMultiLineEntry()
	d.responseBody.Disable()

	d.respBodyArea = container.NewAppTabs(
		container.NewTabItem("Body",
			container.NewBorder(d.responseBodyType, nil, nil, nil, d.responseBody)),
		container.NewTabItem("Header", d.responseHeader),
	)
}

func (d *devToolHttpClient) loadHttpRequest(req *common.HttpRequest) {

}

func filteredListForSelect(filter *string) (resultList []string) {
	allOptionsList := []string{"hello", "world", "123"}
	if filter == nil || *filter == "" {
		return allOptionsList
	} else {
		for _, option := range allOptionsList {
			if strings.Contains(option, *filter) {
				resultList = append(resultList, option)
			}
		}
		return resultList
	}
}
