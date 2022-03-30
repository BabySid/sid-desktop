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

func (d *devToolDateTime) CreateView() fyne.CanvasObject {
	if d.content != nil {
		return d.content
	}

	// common components
	d.controlBtn = widget.NewButton(sidTheme.AppDevToolsDateTimeStartBtnName, func() {

	})

	// from ts to datetime
	d.nowTimeStampBtn = widget.NewButton("1648642283", func() {
		// now timestamp
		// tap it to copy to d.fromTimeStampEntry
	})
	d.nowTimeStampBtn.Alignment = widget.ButtonAlignLeading

	d.timeStampUnitTsToDt = widget.NewSelect([]string{}, func(s string) {

	})
	d.dateTimeFormatTsToDt = widget.NewSelect([]string{}, func(s string) {

	})

	d.fromTimeStampEntry = widget.NewEntry()
	d.toDateTimeEntry = widget.NewEntry()

	fromTsToDtCont := widget.NewForm()
	fromTsToDtCont.Append(sidTheme.AppDevToolsDateTimeNowTSName, d.nowTimeStampBtn)
	fromTsToDtCont.Append(sidTheme.AppDevToolsDateTimeTimeStampName, d.fromTimeStampEntry)
	fromTsToDtCont.Append(sidTheme.AppDevToolsDateTimeDateTimeName, d.toDateTimeEntry)
	fromTsToDtCont.SubmitText = sidTheme.AppDevToolsDateTimeConvertBtnName
	fromTsToDtCont.OnSubmit = func() {
		fmt.Println("fromTsToDtCont")
	}

	fromTsToDtCard := widget.NewCard("", sidTheme.AppDevToolsDateTimeFromTsToDtTitle, fromTsToDtCont)

	// from datetime to ts
	d.nowDateTimeBtn = widget.NewButton("2020-01-10 12:23:23", func() {

	})
	d.nowDateTimeBtn.Alignment = widget.ButtonAlignLeading

	d.timeStampUnitDtToTs = widget.NewSelect([]string{}, func(s string) {

	})
	d.dateTimeFormatDtToTs = widget.NewSelect([]string{}, func(s string) {

	})

	d.fromDateTimeEntry = widget.NewEntry()
	d.toTimeStampEntry = widget.NewEntry()

	fromDtToTsCont := widget.NewForm()
	fromDtToTsCont.Append(sidTheme.AppDevToolsDateTimeNowDateTimeName, d.nowDateTimeBtn)
	fromDtToTsCont.Append(sidTheme.AppDevToolsDateTimeDateTimeName, d.fromDateTimeEntry)
	fromDtToTsCont.Append(sidTheme.AppDevToolsDateTimeTimeStampName, d.toTimeStampEntry)
	fromDtToTsCont.SubmitText = sidTheme.AppDevToolsDateTimeConvertBtnName
	fromDtToTsCont.OnSubmit = func() {
		fmt.Println("fromDtToTsCont")
	}
	fromDtToTsCard := widget.NewCard("", sidTheme.AppDevToolsDateTimeFromDtToTsTitle, fromDtToTsCont)

	d.content = container.NewBorder(container.NewHBox(layout.NewSpacer(), d.controlBtn), nil, nil, nil,
		container.NewGridWithRows(2, fromTsToDtCard, fromDtToTsCard))
	return d.content
}
