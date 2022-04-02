package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"sid-desktop/desktop/common"
	"strings"
)

var _ devToolInterface = (*devToolHttpClient)(nil)

type devToolHttpClient struct {
	method      *widget.Select
	url         *widget.Entry
	sendRequest *widget.Button

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

	d.createRequestView()
	d.createResponseView()

	area := container.NewVSplit(d.reqBodyArea, d.respBodyArea)
	area.SetOffset(0.5)

	d.content = container.NewBorder(
		container.NewBorder(nil, nil, d.method, d.sendRequest, d.url),
		nil, nil, nil,
		area)

	return d.content
}

func (d *devToolHttpClient) createRequestView() {
	d.reqHeaderBinding = binding.NewUntypedList()

	d.reqHeaderBinding.Set(common.NewBuiltInHttpHeaderBinding())
	d.reqHeaderBinding.Append(common.NewHttpHeaderBinding())

	d.requestHeader = widget.NewListWithData(
		d.reqHeaderBinding,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil, widget.NewButton("Remove", nil),
				container.NewGridWithColumns(2,
					widget.NewEntry(),
					widget.NewEntry(),
				))
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			o, _ := item.(binding.Untyped).Get()
			header := o.(*common.HttpHeaderBinding)

			key := obj.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Entry)
			key.Bind(header.Key)
			key.SetPlaceHolder("key")
			key.OnChanged = func(s string) {
				header.Key.Set(s)
				//d.reqHeaderBinding.set
				if s != "" {
					if d.reqHeaderBinding.Length() <= 1 {
						d.reqHeaderBinding.Append(common.NewHttpHeaderBinding())
					}
				}
			}

			value := obj.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Entry)
			value.Bind(header.Value)
			value.SetPlaceHolder("value")

			rm := obj.(*fyne.Container).Objects[1].(*widget.Button)

			rm.Enable()
			if d.reqHeaderBinding.Length() == 1 {
				rm.Disable()
			}

			rm.OnTapped = func() {
				//tmp, _ := d.reqHeaderBinding.Get()
				//for i, h := range tmp {
				//	if h == o {
				//		if i+1 > len(tmp) {
				//			tmp = append(tmp[:i], tmp[:]...)
				//		} else {
				//			tmp = append(tmp[:i], tmp[i+1:]...)
				//		}
				//
				//		d.reqHeaderBinding.Set(tmp)
				//	}
				//}
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
	d.respHeaderBinding = binding.NewUntypedList()
	d.respHeaderBinding.Set(common.NewBuiltInHttpHeaderBinding())

	d.responseHeader = widget.NewListWithData(
		d.reqHeaderBinding,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil, nil,
				container.NewGridWithColumns(2,
					widget.NewEntry(),
					widget.NewEntry(),
				))
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			o, _ := item.(binding.Untyped).Get()
			header := o.(*common.HttpHeaderBinding)

			key := obj.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Entry)
			key.Bind(header.Key)
			key.Disable()

			value := obj.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Entry)
			value.Bind(header.Value)
			value.Disable()
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
