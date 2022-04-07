package gui

import (
	"encoding/base64"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	"sid-desktop/desktop/common"
	sidTheme "sid-desktop/desktop/theme"
	sidWidget "sid-desktop/desktop/widget"
	"strings"
	"time"
)

var _ devToolInterface = (*devToolHttpClient)(nil)

type devToolHttpClient struct {
	devToolAdapter

	method        *widget.Select
	url           *widget.Entry
	sendRequest   *widget.Button
	searchHistory *widget.Button

	// request
	reqBodyArea      *container.AppTabs
	reqHeaderBinding binding.UntypedList
	requestHeader    *widget.List
	requestBody      *widget.Entry
	requestBodyType  *widget.RadioGroup
	prettyReqJson    *widget.Button
	// req-auth
	reqAuthTab           *sidWidget.SelectTab
	reqNoAuthPanel       fyne.CanvasObject
	reqBasicAuthPanel    fyne.CanvasObject
	requestBasicAuthUser *widget.Entry
	requestBasicAuthPass *widget.Entry

	// response
	respBodyArea      *container.AppTabs
	respHeaderBinding binding.UntypedList
	responseHeader    *widget.List
	responseBody      *widget.Entry
	responseBodyType  *widget.RadioGroup
	prettyRespJson    *widget.Button
	respStatus        *widget.Label
}

func (d *devToolHttpClient) CreateView() fyne.CanvasObject {
	if d.content != nil {
		return d.content
	}

	d.method = widget.NewSelect(common.HttpMethod, nil)
	d.method.PlaceHolder = d.method.Options[0]
	d.method.SetSelectedIndex(0)

	d.url = widget.NewEntry()
	d.url.SetPlaceHolder(sidTheme.AppDevToolsHttpCliUrlPlaceHolder)
	d.url.OnChanged = func(s string) {
		if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
			txt := "http://" + s
			d.url.CursorColumn = len(txt)
			d.url.SetText(txt)
		}
	}
	d.sendRequest = widget.NewButtonWithIcon(sidTheme.AppDevToolsHttpCliSendReqName,
		sidTheme.ResourceHttpIcon,
		d.sendHttpRequest)
	d.searchHistory = widget.NewButtonWithIcon(sidTheme.AppDevToolsHttpCliSearchHisName,
		sidTheme.ResourceSearchIcon,
		d.searchFromHistory)

	d.createRequestView()
	d.createResponseView()

	area := container.NewVSplit(d.reqBodyArea, d.respBodyArea)
	area.SetOffset(0.6)

	d.content = container.NewBorder(
		container.NewBorder(nil, nil, d.method, container.NewHBox(d.sendRequest, d.searchHistory), d.url),
		nil, nil, nil,
		area)

	return d.content
}

func (d *devToolHttpClient) createRequestView() {
	d.createRequestHeader()
	d.createRequestBody()
	d.createRequestAuth()

	d.reqBodyArea = container.NewAppTabs(
		container.NewTabItem(sidTheme.AppDevToolsHttpCliBodyTabName,
			container.NewBorder(container.NewHBox(d.requestBodyType, layout.NewSpacer(), d.prettyReqJson),
				nil, nil, nil, d.requestBody)),
		container.NewTabItem(sidTheme.AppDevToolsHttpCliHeaderTabName, d.requestHeader),
		container.NewTabItem(sidTheme.AppDevToolsHttpCliAuthTabName, d.reqAuthTab.Content),
	)
}

func (d *devToolHttpClient) createRequestHeader() {
	d.reqHeaderBinding = binding.NewUntypedList()

	d.reqHeaderBinding.Set(common.NewBuiltInHttpHeader())
	d.reqHeaderBinding.Append(common.NewHttpHeader())

	d.requestHeader = widget.NewListWithData(
		d.reqHeaderBinding,
		func() fyne.CanvasObject {
			key := widget.NewSelectEntry(common.BuiltInHttpHeaderName())
			key.SetPlaceHolder(sidTheme.AppDevToolsHttpCliReqHeaderKeyPlaceHolder)
			value := widget.NewEntry()
			value.SetPlaceHolder(sidTheme.AppDevToolsHttpCliReqHeaderValPlaceHolder)
			return container.NewBorder(nil, nil, nil,
				widget.NewButtonWithIcon(sidTheme.AppDevToolsHttpCliRmReqHeaderName, sidTheme.ResourceRmIcon, nil),
				container.NewGridWithColumns(2,
					key,
					value,
				))
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			o, _ := item.(binding.Untyped).Get()
			header := o.(*common.HttpHeader)

			arr, _ := d.reqHeaderBinding.Get()
			lineNo := gobase.ContainsInterface(arr, header)
			gobase.True(lineNo >= 0)

			key := obj.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.SelectEntry)
			// OnChanged must set nil because SetText will call OnChanged if the method is not nil.
			// This can cause display confusion
			// Another way is to use key.Bind(), but it is very inconvenient.
			// On the one hand, it needs to call Unbind(), and on the other hand, it also needs to call RemoveListener().
			// More importantly, the function of the automatic add row is expected to be called after SetText.
			key.OnChanged = nil
			key.SetText(header.Key)
			key.OnChanged = func(s string) {
				header.Key = s
				if lineNo == d.reqHeaderBinding.Length()-1 {
					_ = d.reqHeaderBinding.Append(common.NewHttpHeader())
				}
				key.SetOptions(common.FilterOption(s, common.BuiltInHttpHeaderName()))
			}

			value := obj.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Entry)
			value.OnChanged = nil
			value.SetText(header.Value)
			value.OnChanged = func(s string) {
				header.Value = s
			}

			rm := obj.(*fyne.Container).Objects[1].(*widget.Button)
			rm.Enable()
			if d.reqHeaderBinding.Length() == 1 {
				rm.Disable()
			}
			rm.OnTapped = func() {
				tmp, _ := d.reqHeaderBinding.Get()
				tmp = append(tmp[:lineNo], tmp[lineNo+1:]...)
				_ = d.reqHeaderBinding.Set(tmp)
			}
		},
	)
}

func (d *devToolHttpClient) createRequestBody() {
	d.prettyReqJson = widget.NewButtonWithIcon(sidTheme.AppDevToolsCliPrettyJsonName, sidTheme.ResourcePrettyIcon, func() {
		if d.requestBody.Text == "" {
			return
		}

		out, err := gobase.PrettyPrintJson(d.requestBody.Text, prettyJsonIndent)
		if err != nil {
			dialog.ShowError(err, globalWin.win)
			return
		}
		d.requestBody.SetText(out)
	})

	d.requestBodyType = widget.NewRadioGroup([]string{
		sidTheme.AppDevToolsHttpCliBodyTypeName1,
		sidTheme.AppDevToolsHttpCliBodyTypeName2,
	}, func(s string) {
		if s == sidTheme.AppDevToolsHttpCliBodyTypeName2 {
			d.prettyReqJson.Enable()
		} else {
			d.prettyReqJson.Disable()
		}
	})
	d.requestBodyType.Horizontal = true
	d.requestBodyType.Required = true
	d.requestBodyType.SetSelected(sidTheme.AppDevToolsHttpCliBodyTypeName2)

	d.requestBody = widget.NewMultiLineEntry()
}

func (d *devToolHttpClient) createRequestAuth() {
	d.reqNoAuthPanel = container.NewBorder(nil, nil, nil, nil, layout.NewSpacer())
	d.requestBasicAuthUser = widget.NewEntry()
	d.requestBasicAuthUser.OnChanged = func(s string) {
		d.updateBasicAuthForReqHeader()
	}
	d.requestBasicAuthPass = widget.NewPasswordEntry()
	d.requestBasicAuthPass.OnChanged = func(s string) {
		d.updateBasicAuthForReqHeader()
	}
	d.reqBasicAuthPanel = widget.NewForm(
		widget.NewFormItem(sidTheme.AppDevToolsHttpCliBasicAuthUser, d.requestBasicAuthUser),
		widget.NewFormItem(sidTheme.AppDevToolsHttpCliBasicAuthPass, d.requestBasicAuthPass),
	)
	d.reqAuthTab = sidWidget.NewSelectTab(
		sidWidget.NewSelectItem(sidTheme.AppDevToolsHttpCliAuthTypeName1, d.reqNoAuthPanel),
		sidWidget.NewSelectItem(sidTheme.AppDevToolsHttpCliAuthTypeName2, d.reqBasicAuthPanel),
	)
	d.reqAuthTab.OnSelected = func(item *sidWidget.SelectItem) {
		switch item.Text {
		case sidTheme.AppDevToolsHttpCliAuthTypeName2:
			d.addBasicAuthForReqHeader()
		default:
			d.rmBasicAuthForReqHeader()
		}
	}
}

func (d *devToolHttpClient) updateBasicAuthForReqHeader() {
	line, header := d.findAuthHeaderForReqHeader()
	gobase.True(line >= 0 && header != nil)

	header.Value = "Basic " + base64.StdEncoding.EncodeToString([]byte(d.requestBasicAuthUser.Text+":"+d.requestBasicAuthPass.Text))
}

func (d *devToolHttpClient) addBasicAuthForReqHeader() {
	// user may add it manually
	line, header := d.findAuthHeaderForReqHeader()
	if line >= 0 || header != nil {
		return
	}

	auth := common.AuthHeader
	auth.Value = "Basic " + base64.StdEncoding.EncodeToString([]byte(d.requestBasicAuthUser.Text+":"+d.requestBasicAuthPass.Text))
	_ = d.reqHeaderBinding.Prepend(&auth)
}

func (d *devToolHttpClient) rmBasicAuthForReqHeader() {
	// user may remove it manually
	idx, header := d.findAuthHeaderForReqHeader()
	if idx < 0 || header == nil {
		return
	}

	headers, _ := d.reqHeaderBinding.Get()
	headers = append(headers[:idx], headers[idx+1:]...)
	_ = d.reqHeaderBinding.Set(headers)
}

func (d *devToolHttpClient) findAuthHeaderForReqHeader() (int, *common.HttpHeader) {
	headers, _ := d.reqHeaderBinding.Get()

	for i, obj := range headers {
		header := obj.(*common.HttpHeader)
		if header.Key == common.AuthHeader.Key {
			return i, header
		}
	}

	return -1, nil
}

func (d *devToolHttpClient) createResponseView() {
	d.respHeaderBinding = binding.NewUntypedList()

	d.responseHeader = widget.NewListWithData(
		d.respHeaderBinding,
		func() fyne.CanvasObject {
			k := widget.NewEntry()
			v := widget.NewEntry()
			k.Disable()
			v.Disable()
			return container.NewBorder(nil, nil, nil, nil,
				container.NewGridWithColumns(2,
					k,
					v,
				))
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			o, _ := item.(binding.Untyped).Get()
			header := o.(*common.HttpHeader)

			key := obj.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Entry)
			key.SetText(header.Key)

			value := obj.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Entry)
			value.SetText(header.Value)
		},
	)

	d.prettyRespJson = widget.NewButtonWithIcon(sidTheme.AppDevToolsCliPrettyJsonName, sidTheme.ResourcePrettyIcon, func() {
		if d.responseBody.Text == "" {
			return
		}

		out, err := gobase.PrettyPrintJson(d.responseBody.Text, prettyJsonIndent)
		if err != nil {
			dialog.ShowError(err, globalWin.win)
			return
		}
		d.responseBody.SetText(out)
	})

	d.responseBodyType = widget.NewRadioGroup([]string{
		sidTheme.AppDevToolsHttpCliBodyTypeName1,
		sidTheme.AppDevToolsHttpCliBodyTypeName2,
	}, nil)
	d.responseBodyType.Horizontal = true
	d.responseBodyType.Required = true
	d.responseBodyType.SetSelected(sidTheme.AppDevToolsHttpCliBodyTypeName2)
	d.responseBodyType.Disable()

	d.respStatus = widget.NewLabel("")

	d.responseBody = widget.NewMultiLineEntry()
	d.responseBody.Disable()

	d.respBodyArea = container.NewAppTabs(
		container.NewTabItem(sidTheme.AppDevToolsHttpCliBodyTabName,
			container.NewBorder(
				container.NewHBox(d.responseBodyType, layout.NewSpacer(), d.respStatus, d.prettyRespJson), nil, nil, nil,
				d.responseBody)),
		container.NewTabItem(sidTheme.AppDevToolsHttpCliHeaderTabName, d.responseHeader),
	)
}

func (d *devToolHttpClient) sendHttpRequest() {
	header := make(map[string]string)
	arr, _ := d.reqHeaderBinding.Get()
	for _, item := range arr {
		h := item.(*common.HttpHeader)
		header[h.Key] = h.Value
	}

	begin := time.Now()
	status, respHeader, body, err := common.DoHttpRequest(d.method.Selected, d.url.Text, d.requestBody.Text, header)
	cost := time.Since(begin)
	if err != nil {
		d.responseBody.SetText(err.Error())
		return
	}
	d.respStatus.SetText(fmt.Sprintf(sidTheme.AppDevToolsCliRespStatusFormat, status, cost, len(body)))

	d.responseBodyType.SetSelected(sidTheme.AppDevToolsHttpCliBodyTypeName1)
	d.responseBody.SetText(string(body))
	rs := make([]interface{}, 0)
	for k, v := range respHeader {
		value := strings.Join(v, " ")
		if k == "Content-Type" && strings.Index(value, "application/json") >= 0 {
			d.responseBodyType.SetSelected(sidTheme.AppDevToolsHttpCliBodyTypeName2)
		}
		header := &common.HttpHeader{
			Key:   k,
			Value: value,
		}
		rs = append(rs, header)
	}
	d.respHeaderBinding.Set(rs)
}

func (d *devToolHttpClient) searchFromHistory() {

}

func (d *devToolHttpClient) loadHttpRequest(req *common.HttpRequest) {

}
