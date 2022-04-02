package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	"sid-desktop/desktop/common"
	"strconv"
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
	d.sendRequest = widget.NewButton("Send", d.sendHttpRequest)

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

func (d *devToolHttpClient) sendHttpRequest() {
	//for i := 0; i < d.reqHeaderBinding.Length(); i++ {
	//	obj, _ := d.reqHeaderBinding.GetValue(i)
	//	header := obj.(*common.HttpHeader)
	//	fmt.Println(i, header.Key, header.Value)
	//}
	i, _ := strconv.Atoi(d.url.Text)
	arr, _ := d.reqHeaderBinding.Get()
	if len(arr) == 0 {
		return
	}
	arr = append(arr[:i], arr[i+1:]...)
	d.reqHeaderBinding.Set(arr)
}

func (d *devToolHttpClient) createRequestView() {
	d.reqHeaderBinding = binding.NewUntypedList()

	d.reqHeaderBinding.Set(common.NewBuiltInHttpHeader())
	//d.reqHeaderBinding.Append(common.NewHttpHeader())

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
			header := o.(*common.HttpHeader)

			arr, _ := d.reqHeaderBinding.Get()
			lineNo := ContainsInterface(arr, header)
			gobase.True(lineNo >= 0)

			key := obj.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Entry)
			key.SetText(header.Key)
			key.SetPlaceHolder("key")
			//key.OnChanged = func(s string) {
			//	fmt.Printf("set key %d %s->%s\n", lineNo, header.Key, s)
			//
			//	header.Key = s
			//
			//	//if lineNo == d.reqHeaderBinding.Length()-1 {
			//	//	d.reqHeaderBinding.Append(common.NewHttpHeader())
			//	//}
			//}

			value := obj.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Entry)
			if header.Value != nil {
				value.SetText(fmt.Sprintf("%v", header.Value))
			} else {
				value.SetText("")
			}

			value.SetPlaceHolder("value")
			//value.OnChanged = func(s string) {
			//	fmt.Printf("set value %d %s->%s\n", lineNo, header.Value, s)
			//	header.Value = s
			//}

			rm := obj.(*fyne.Container).Objects[1].(*widget.Button)

			rm.Enable()
			if d.reqHeaderBinding.Length() == 1 {
				rm.Disable()
			}

			rm.OnTapped = func() {
				tmp, _ := d.reqHeaderBinding.Get()
				tmp = append(tmp[:lineNo], tmp[lineNo+1:]...)
				d.reqHeaderBinding.Set(tmp)
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
	d.respHeaderBinding.Set(common.NewBuiltInHttpHeader())

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
			header := o.(*common.HttpHeader)

			key := obj.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Entry)
			key.SetText(header.Key)
			key.Disable()

			value := obj.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Entry)
			value.SetText(fmt.Sprintf("%v", header.Value))
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

func ContainsInterface(array []interface{}, val interface{}) int {
	for i, item := range array {
		if item == val {
			return i
		}
	}

	return -1
}
