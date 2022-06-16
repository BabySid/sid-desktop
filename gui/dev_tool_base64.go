package gui

import (
	"encoding/base64"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	"sid-desktop/theme"
)

var _ devToolInterface = (*devToolBase64)(nil)

type devToolBase64 struct {
	devToolAdapter
	urlChk  *widget.Check
	fillChk *widget.Check

	encodeBtn *widget.Button
	decodeBtn *widget.Button

	inputText  *widget.Entry
	outputText *widget.Entry
}

func (d *devToolBase64) CreateView() fyne.CanvasObject {
	if d.content != nil {
		return d.content
	}

	d.encodeBtn = widget.NewButtonWithIcon(theme.AppDevToolsBase64EncodeName, theme.ResourceEncodeIcon, func() {
		enc := d.getEncoding()
		out := enc.EncodeToString([]byte(d.inputText.Text))
		d.outputText.SetText(out)
	})
	d.decodeBtn = widget.NewButtonWithIcon(theme.AppDevToolsBase64DecodeName, theme.ResourceDecodeIcon, func() {
		enc := d.getEncoding()

		out, err := enc.DecodeString(d.inputText.Text)
		if err != nil {
			d.outputText.SetText(err.Error())
		} else {
			d.outputText.SetText(string(out))
		}
	})

	d.urlChk = widget.NewCheck(theme.AppDevToolsBase64UrlModeName, nil)
	d.urlChk.SetChecked(false)

	d.fillChk = widget.NewCheck(theme.AppDevToolsBase64FillModeName, nil)
	d.fillChk.SetChecked(true)

	d.inputText = widget.NewMultiLineEntry()
	d.inputText.Wrapping = fyne.TextWrapWord
	d.inputText.SetPlaceHolder(theme.AppDevToolsBase64InputName)

	d.outputText = widget.NewMultiLineEntry()
	d.outputText.Wrapping = fyne.TextWrapWord
	d.outputText.SetPlaceHolder(theme.AppDevToolsBase64OutputName)

	cont := container.NewHSplit(d.inputText, d.outputText)
	cont.SetOffset(0.5)
	d.content = container.NewBorder(container.NewHBox(d.encodeBtn, d.decodeBtn, layout.NewSpacer(), d.urlChk, d.fillChk),
		nil, nil, nil, cont)

	return d.content
}

func (d *devToolBase64) getEncoding() *base64.Encoding {
	if d.urlChk.Checked && d.fillChk.Checked {
		return base64.URLEncoding
	}
	if d.urlChk.Checked && !d.fillChk.Checked {
		return base64.RawURLEncoding
	}
	if !d.urlChk.Checked && d.fillChk.Checked {
		return base64.StdEncoding
	}
	if !d.urlChk.Checked && !d.fillChk.Checked {
		return base64.RawStdEncoding
	}

	gobase.AssertHere()
	return nil
}
