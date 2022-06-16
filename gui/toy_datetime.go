package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	sidTheme "sid-desktop/theme"
	"time"
)

var _ toyInterface = (*toyDateTime)(nil)

type toyDateTime struct {
	toyAdapter
	date *canvas.Text
	week *canvas.Text
	time *canvas.Text
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

	_ = gobase.GlobalScheduler.AddJob("toy_datetime", "*/1 * * * * *", tdt)
	// todo assert(err == nil)
	return nil
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

	tdt.date.Text = gobase.FormatTimeStampWithFormat(now.Unix(), gobase.DateFormat)
	tdt.date.Refresh()

	tdt.week.Text = weekStr[now.Weekday()]
	tdt.week.Refresh()

	tdt.time.Text = gobase.FormatTimeStampWithFormat(now.Unix(), gobase.TimeFormat)
	tdt.time.Refresh()
}
