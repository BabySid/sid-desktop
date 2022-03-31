package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var _ devToolInterface = (*devToolHttpClient)(nil)

type devToolHttpClient struct {
	method      *widget.Select
	url         *widget.Entry
	sendRequest *widget.Button

	requestHeader   *widget.List
	requestBody     *widget.Entry
	requestBodyType *widget.RadioGroup

	responseHeader   *widget.Table
	responseBody     *widget.Entry
	responseBodyType *widget.RadioGroup

	content fyne.CanvasObject
}

func (d *devToolHttpClient) CreateView() fyne.CanvasObject {
	if d.content != nil {
		return d.content
	}

	d.requestHeader = widget.NewList(
		func() int {
			return 1
		},
		func() fyne.CanvasObject {
			return container.NewGridWithColumns(5,
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				container.NewHBox(layout.NewSpacer(), widget.NewSeparator(), layout.NewSpacer()),
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				container.NewHBox(layout.NewSpacer(), widget.NewSeparator(), layout.NewSpacer()),
				widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{}))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText("sdasdad")
			item.(*fyne.Container).Objects[2].(*widget.Label).SetText("xxx")
			item.(*fyne.Container).Objects[4].(*widget.Label).SetText("adsadas")
		},
	)

	d.content = container.NewBorder(nil, nil, nil, nil,
		d.requestHeader)
	return d.content
}
