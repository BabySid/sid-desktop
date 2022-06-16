package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	"sid-desktop/theme"
	"strconv"
	"time"
)

var _ devToolInterface = (*devToolDateTime)(nil)

type devToolDateTime struct {
	devToolAdapter
	// common components
	controlBtn *widget.Button
	done       chan bool
	timer      *time.Ticker

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
	d.controlBtn = widget.NewButtonWithIcon(theme.AppDevToolsDateTimeStopBtnName, theme.ResourceStopIcon, func() {
		if d.controlBtn.Text == theme.AppDevToolsDateTimeStopBtnName {
			d.stopTimer()
			d.controlBtn.SetText(theme.AppDevToolsDateTimeStartBtnName)
			d.controlBtn.SetIcon(theme.ResourceRunIcon)
		} else {
			d.startTimer()
			d.controlBtn.SetText(theme.AppDevToolsDateTimeStopBtnName)
			d.controlBtn.SetIcon(theme.ResourceStopIcon)
		}
	})

	// from ts to datetime
	d.nowTimeStampBtn = widget.NewButton("", func() {
		d.fromTimeStampEntry.SetText(d.nowTimeStampBtn.Text)
	})
	d.nowTimeStampBtn.Alignment = widget.ButtonAlignLeading

	d.dateTimeFormatTsToDt = widget.NewSelect([]string{}, nil)
	d.timeStampUnitTsToDt = widget.NewSelect(timeStampUnit, func(s string) {
		if s == tsUnit1 {
			d.dateTimeFormatTsToDt.Options = dateTimeFormatS
		}
		if s == tsUnit2 {
			d.dateTimeFormatTsToDt.Options = dateTimeFormatMS
		}
		setSelectLayout(d.dateTimeFormatTsToDt)
		d.dateTimeFormatTsToDt.SetSelectedIndex(0)
	})
	setSelectLayout(d.timeStampUnitTsToDt)
	d.timeStampUnitTsToDt.SetSelectedIndex(0)

	d.fromTimeStampEntry = widget.NewEntry()
	d.toDateTimeEntry = widget.NewEntry()

	fromTsToDtCont := widget.NewForm(
		widget.NewFormItem(theme.AppDevToolsDateTimeTSUnit, d.timeStampUnitTsToDt),
		widget.NewFormItem(theme.AppDevToolsDateTimeNowTSName, d.nowTimeStampBtn),
		widget.NewFormItem(theme.AppDevToolsDateTimeTimeStampName, d.fromTimeStampEntry),
		widget.NewFormItem(theme.AppDevToolsDateTimeDTFormat, d.dateTimeFormatTsToDt),
		widget.NewFormItem(theme.AppDevToolsDateTimeDateTimeName, d.toDateTimeEntry),
	)

	fromTsToDtCont.SubmitText = theme.AppDevToolsDateTimeConvertBtnName
	fromTsToDtCont.OnSubmit = func() {
		ts, err := strconv.ParseInt(d.fromTimeStampEntry.Text, 10, 64)
		if err != nil {
			d.toDateTimeEntry.SetText(err.Error())
			return
		}

		if d.timeStampUnitTsToDt.Selected == tsUnit1 {
			d.toDateTimeEntry.SetText(gobase.FormatTimeStampWithFormat(ts, d.dateTimeFormatTsToDt.Selected))
		}
		if d.timeStampUnitTsToDt.Selected == tsUnit2 {
			d.toDateTimeEntry.SetText(gobase.FormatTimeStampMilliWithFormat(ts, d.dateTimeFormatTsToDt.Selected))
		}
	}

	fromTsToDtCard := widget.NewCard("", theme.AppDevToolsDateTimeFromTsToDtTitle, fromTsToDtCont)

	// from datetime to ts
	d.nowDateTimeBtn = widget.NewButton("", func() {
		d.fromDateTimeEntry.SetText(d.nowDateTimeBtn.Text)
	})
	d.nowDateTimeBtn.Alignment = widget.ButtonAlignLeading

	d.dateTimeFormatDtToTs = widget.NewSelect([]string{}, nil)

	d.timeStampUnitDtToTs = widget.NewSelect(timeStampUnit, func(s string) {
		if s == tsUnit1 {
			d.dateTimeFormatDtToTs.Options = dateTimeFormatS
		}
		if s == tsUnit2 {
			d.dateTimeFormatDtToTs.Options = dateTimeFormatMS
		}
		setSelectLayout(d.dateTimeFormatDtToTs)
		d.dateTimeFormatDtToTs.SetSelectedIndex(0)
	})
	setSelectLayout(d.timeStampUnitDtToTs)
	d.timeStampUnitDtToTs.SetSelectedIndex(0)

	d.fromDateTimeEntry = widget.NewEntry()
	d.toTimeStampEntry = widget.NewEntry()

	fromDtToTsCont := widget.NewForm(
		widget.NewFormItem(theme.AppDevToolsDateTimeDTFormat, d.dateTimeFormatDtToTs),
		widget.NewFormItem(theme.AppDevToolsDateTimeNowDateTimeName, d.nowDateTimeBtn),
		widget.NewFormItem(theme.AppDevToolsDateTimeDateTimeName, d.fromDateTimeEntry),
		widget.NewFormItem(theme.AppDevToolsDateTimeTSUnit, d.timeStampUnitDtToTs),
		widget.NewFormItem(theme.AppDevToolsDateTimeTimeStampName, d.toTimeStampEntry),
	)
	fromDtToTsCont.SubmitText = theme.AppDevToolsDateTimeConvertBtnName
	fromDtToTsCont.OnSubmit = func() {
		t, err := time.Parse(d.dateTimeFormatDtToTs.Selected, d.fromDateTimeEntry.Text)
		if err != nil {
			d.toTimeStampEntry.SetText(err.Error())
			return
		}

		if d.timeStampUnitDtToTs.Selected == tsUnit1 {
			d.toTimeStampEntry.SetText(strconv.Itoa(int(t.Unix())))
		}
		if d.timeStampUnitDtToTs.Selected == tsUnit2 {
			d.toTimeStampEntry.SetText(strconv.FormatInt(t.UnixMilli(), 10))
		}
	}

	fromDtToTsCard := widget.NewCard("", theme.AppDevToolsDateTimeFromDtToTsTitle, fromDtToTsCont)

	d.content = container.NewBorder(container.NewHBox(layout.NewSpacer(), d.controlBtn), nil, nil, nil,
		container.NewScroll(container.NewVBox(fromTsToDtCard, fromDtToTsCard, layout.NewSpacer())))

	d.done = make(chan bool)
	d.startTimer()

	return d.content
}

// getMaxLengthItem returns the maxOption for widget.Select Layout
// https://github.com/fyne-io/fyne/issues/2881
func setSelectLayout(sel *widget.Select) {
	value := ""
	for _, item := range sel.Options {
		if len(item) > len(value) {
			value = item
		}
	}

	sel.PlaceHolder = value
}

func (d *devToolDateTime) startTimer() {
	d.timer = time.NewTicker(time.Second * 1)
	go func() {
		defer func() {
			d.timer.Stop()
		}()

		for {
			select {
			case <-d.timer.C:
				now := time.Now()
				if d.timeStampUnitTsToDt.Selected == tsUnit1 {
					d.nowTimeStampBtn.SetText(strconv.Itoa(int(now.Unix())))
				}
				if d.timeStampUnitTsToDt.Selected == tsUnit2 {
					d.nowTimeStampBtn.SetText(strconv.Itoa(int(now.UnixMilli())))
				}

				if d.timeStampUnitDtToTs.Selected == tsUnit1 {
					d.nowDateTimeBtn.SetText(gobase.FormatTimeStampWithFormat(now.Unix(), d.dateTimeFormatDtToTs.Selected))
				}
				if d.timeStampUnitDtToTs.Selected == tsUnit2 {
					d.nowDateTimeBtn.SetText(gobase.FormatTimeStampMilliWithFormat(now.UnixMilli(), d.dateTimeFormatDtToTs.Selected))
				}

			case <-d.done:
				return
			}
		}
	}()
}

func (d *devToolDateTime) stopTimer() {
	d.done <- true
}
