package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"sid-desktop/base"
	sidTheme "sid-desktop/desktop/theme"
	"time"
)

var _ toyInterface = (*toyDateTime)(nil)

type toyDateTime struct {
	date *canvas.Text
	week *canvas.Text
	time *canvas.Text

	widget *widget.Card
}

func (tdt *toyDateTime) Init() error {
	tdt.date = canvas.NewText("", nil)
	tdt.date.TextSize = 15

	tdt.week = canvas.NewText("", nil)
	tdt.week.TextSize = 15

	tdt.time = canvas.NewText("", nil)
	tdt.time.TextSize = 42

	tdt.widget = widget.NewCard("", sidTheme.ToyDateTimeTitle,
		container.NewVBox(
			container.NewHBox(layout.NewSpacer(), tdt.date, tdt.week, layout.NewSpacer()),
			container.NewCenter(tdt.time),
		),
	)

	tdt.widget.Resize(fyne.NewSize(ToyWidth, 130))

	// for init
	tdt.Run()

	_ = base.GlobalScheduler.AddJob("toy_datetime", "*/1 * * * * *", tdt)
	// todo assert(err == nil)
	return nil
}

func (tdt *toyDateTime) GetToyCard() *widget.Card {
	return tdt.widget
}

var (
	weekStr = [...]string{
		"Sunday",
		"Monday",
		"Tuesday",
		"Wednesday",
		"Thursday",
		"Friday",
		"Saturday",
	}
)

func (tdt *toyDateTime) Run() {
	now := time.Now()

	tdt.date.Text = base.FormatTimeStampWithFormat(now.Unix(), base.DateFormat)
	tdt.date.Refresh()

	tdt.week.Text = weekStr[now.Weekday()]
	tdt.week.Refresh()

	tdt.time.Text = base.FormatTimeStampWithFormat(now.Unix(), base.TimeFormat)
	tdt.time.Refresh()
}
