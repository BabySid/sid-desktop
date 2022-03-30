package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	sidTheme "sid-desktop/desktop/theme"
)

var _ devToolInterface = (*devToolDateTime)(nil)

type devToolDateTime struct {
	// common components
	controlBtn *widget.Button

	// from ts to datetime
	nowTimeStampBtn *widget.Button

	timeStampUnitTsToDt  *widget.Select
	dateTimeFormatTsToDt *widget.Select

	fromTimeStampEntry *widget.Entry
	toDateTimeEntry    *widget.Entry

	// from datetime to ts
	nowDateTimeBtn *widget.Button

	timeStampUnitDtToTs  *widget.Select
	dateTimeFormatDtToTs *widget.Select

	fromDateTimeEntry *widget.Entry
	toTimeStampEntry  *widget.Entry

	content fyne.CanvasObject
}

const (
	tsUnit1 = "Second"
	tsUnit2 = "MilliSecond"

	dtLayoutS1 = "2006-01-02 15:04:05"
	dtLayoutS2 = "2006/01/02 15:04:05"
	dtLayoutS3 = "2006-01-02"
	dtLayoutS4 = "2006/01/02"
	dtLayoutS5 = "15:04:05"

	dtLayoutMS1 = "2006-01-02 15:04:05.000"
	dtLayoutMS2 = "2006/01/02 15:04:05.000"
	dtLayoutMS3 = "2006-01-02"
	dtLayoutMS4 = "2006/01/02"
	dtLayoutMS5 = "15:04:05.000"
)

var (
	timeStampUnit = []string{
		tsUnit1,
		tsUnit2,
	}
	dateTimeFormatS = []string{
		dtLayoutS1,
		dtLayoutS2,
		dtLayoutS3,
		dtLayoutS4,
		dtLayoutS5,
	}
	dateTimeFormatMS = []string{
		dtLayoutMS1,
		dtLayoutMS2,
		dtLayoutMS3,
		dtLayoutMS4,
		dtLayoutMS5,
	}
)

func (d *devToolDateTime) CreateView() fyne.CanvasObject {
	if d.content != nil {
		return d.content
	}

	// common components
	d.controlBtn = widget.NewButtonWithIcon(sidTheme.AppDevToolsDateTimeStartBtnName, sidTheme.ResourceRunIcon, func() {

	})

	// from ts to datetime
	d.nowTimeStampBtn = widget.NewButton("1648642283", func() {
		// now timestamp
		// tap it to copy to d.fromTimeStampEntry
	})
	d.nowTimeStampBtn.Alignment = widget.ButtonAlignLeading

	d.timeStampUnitTsToDt = widget.NewSelect(timeStampUnit, func(s string) {
		if s == tsUnit1 {
			d.dateTimeFormatTsToDt.Options = dateTimeFormatS
		}
		if s == tsUnit2 {
			d.dateTimeFormatTsToDt.Options = dateTimeFormatMS
			d.dateTimeFormatTsToDt.PlaceHolder = dtLayoutMS1
		}
		d.dateTimeFormatTsToDt.SetSelectedIndex(0)
	})

	d.dateTimeFormatTsToDt = widget.NewSelect(dateTimeFormatS, nil)

	d.timeStampUnitTsToDt.SetSelectedIndex(0)
	d.dateTimeFormatTsToDt.PlaceHolder = dtLayoutS1
	d.dateTimeFormatTsToDt.SetSelectedIndex(0)

	d.fromTimeStampEntry = widget.NewEntry()
	d.toDateTimeEntry = widget.NewEntry()

	fromTsToDtCont := widget.NewForm(
		widget.NewFormItem(sidTheme.AppDevToolsDateTimeNowTSName, d.nowTimeStampBtn),
		widget.NewFormItem(sidTheme.AppDevToolsDateTimeTimeStampName, d.fromTimeStampEntry),
		widget.NewFormItem(sidTheme.AppDevToolsDateTimeDateTimeName, d.toDateTimeEntry),
	)

	fromTsToDtCont.SubmitText = sidTheme.AppDevToolsDateTimeConvertBtnName
	fromTsToDtCont.OnSubmit = func() {
		fmt.Println("fromTsToDtCont")
	}

	fromTsToDtCard := container.NewVBox(
		container.NewHBox(layout.NewSpacer(), d.timeStampUnitTsToDt, d.dateTimeFormatTsToDt),
		widget.NewCard("", sidTheme.AppDevToolsDateTimeFromTsToDtTitle, fromTsToDtCont),
	)

	// from datetime to ts
	d.nowDateTimeBtn = widget.NewButton("2020-01-10 12:23:23", func() {

	})
	d.nowDateTimeBtn.Alignment = widget.ButtonAlignLeading

	d.timeStampUnitDtToTs = widget.NewSelect(timeStampUnit, func(s string) {
		if s == tsUnit1 {
			d.dateTimeFormatDtToTs.Options = dateTimeFormatS
		}
		if s == tsUnit2 {
			d.dateTimeFormatDtToTs.Options = dateTimeFormatMS
			d.dateTimeFormatDtToTs.PlaceHolder = dtLayoutMS1
		}
		d.dateTimeFormatDtToTs.SetSelectedIndex(0)
	})

	d.dateTimeFormatDtToTs = widget.NewSelect(dateTimeFormatS, nil)

	d.timeStampUnitDtToTs.SetSelectedIndex(0)
	d.dateTimeFormatDtToTs.PlaceHolder = dtLayoutS1
	d.dateTimeFormatDtToTs.SetSelectedIndex(0)

	d.fromDateTimeEntry = widget.NewEntry()
	d.toTimeStampEntry = widget.NewEntry()

	fromDtToTsCont := widget.NewForm(
		widget.NewFormItem(sidTheme.AppDevToolsDateTimeNowDateTimeName, d.nowDateTimeBtn),
		widget.NewFormItem(sidTheme.AppDevToolsDateTimeDateTimeName, d.fromDateTimeEntry),
		widget.NewFormItem(sidTheme.AppDevToolsDateTimeTimeStampName, d.toTimeStampEntry),
	)
	fromDtToTsCont.SubmitText = sidTheme.AppDevToolsDateTimeConvertBtnName
	fromDtToTsCont.OnSubmit = func() {
		fmt.Println("fromDtToTsCont")
	}

	fromDtToTsCard := container.NewVBox(
		container.NewHBox(layout.NewSpacer(), d.timeStampUnitDtToTs, d.dateTimeFormatDtToTs),
		widget.NewCard("", sidTheme.AppDevToolsDateTimeFromDtToTsTitle, fromDtToTsCont),
	)

	d.content = container.NewBorder(container.NewHBox(layout.NewSpacer(), d.controlBtn), nil, nil, nil,
		container.NewVBox(fromTsToDtCard, fromDtToTsCard, layout.NewSpacer()))
	return d.content
}
