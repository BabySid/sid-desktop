package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	sidTheme "sid-desktop/desktop/theme"
)

var _ devToolInterface = (*devToolJson)(nil)

type devToolJson struct {
	devToolAdapter
	compressBtn     *widget.Button
	compressJsonPos *widget.Label
	prettyBtn       *widget.Button
	prettyJsonPos   *widget.Label

	compressJsonText *widget.Entry
	prettyJsonText   *widget.Entry
}

const (
	prettyJsonIndent = "\t"
)

func (d *devToolJson) CreateView() fyne.CanvasObject {
	if d.content != nil {
		return d.content
	}

	d.compressBtn = widget.NewButtonWithIcon(sidTheme.AppDevToolsCompressJsonName, sidTheme.ResourceCompressIcon, func() {
		txt, err := gobase.CompressJson(d.compressJsonText.Text)
		if err != nil {
			d.compressJsonText.SetText(err.Error())
		} else {
			d.compressJsonText.SetText(txt)
		}
	})
	d.prettyBtn = widget.NewButtonWithIcon(sidTheme.AppDevToolsPrettyJsonName, sidTheme.ResourcePrettyIcon, func() {
		txt, err := gobase.PrettyPrintJson(d.prettyJsonText.Text, prettyJsonIndent)
		if err != nil {
			d.prettyJsonText.SetText(err.Error())
		} else {
			d.prettyJsonText.SetText(txt)
		}
	})

	d.compressJsonPos = widget.NewLabel("")
	d.prettyJsonPos = widget.NewLabel("")

	d.compressJsonText = widget.NewMultiLineEntry()
	d.compressJsonText.OnCursorChanged = func() {
		d.compressJsonPos.SetText(
			fmt.Sprintf(sidTheme.TextCursorPosFormat, d.compressJsonText.CursorRow+1, d.compressJsonText.CursorColumn+1))
	}
	d.prettyJsonText = widget.NewMultiLineEntry()
	d.prettyJsonText.OnCursorChanged = func() {
		d.prettyJsonPos.SetText(
			fmt.Sprintf(sidTheme.TextCursorPosFormat, d.prettyJsonText.CursorRow+1, d.prettyJsonText.CursorColumn+1))
	}

	left := container.NewBorder(
		container.NewHBox(d.compressBtn, layout.NewSpacer(), d.compressJsonPos),
		nil, nil, nil, d.compressJsonText)
	right := container.NewBorder(
		container.NewHBox(d.prettyBtn, layout.NewSpacer(), d.prettyJsonPos),
		nil, nil, nil, d.prettyJsonText)

	cont := container.NewHSplit(left, right)
	cont.SetOffset(0.5)
	d.content = cont
	return d.content
}
