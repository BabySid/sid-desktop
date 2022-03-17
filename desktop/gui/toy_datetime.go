package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"sid/base"
	sidTheme "sid/desktop/theme"
	"time"
)
import xwidget "fyne.io/x/fyne/widget"

var _ toyInterface = (*toyDateTime)(nil)

type toyDateTime struct {
	h1  *xwidget.HexWidget
	h2  *xwidget.HexWidget
	dot *widget.RichText
	m1  *xwidget.HexWidget
	m2  *xwidget.HexWidget

	dt *widget.Label

	widget *widget.Card
}

func (tdt *toyDateTime) Init() error {
	tdt.h1 = xwidget.NewHexWidget()
	tdt.h2 = xwidget.NewHexWidget()
	tdt.dot = widget.NewRichTextFromMarkdown("# :")
	tdt.m1 = xwidget.NewHexWidget()
	tdt.m2 = xwidget.NewHexWidget()
	tdt.dt = widget.NewLabel("")

	hexes := []*xwidget.HexWidget{tdt.h1, tdt.h2, tdt.m1, tdt.m2}
	for _, w := range hexes {
		w.SetSlant(0)

		on := color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 0xff,
		}
		w.SetOnColor(on)

		off := color.RGBA{
			R: 1,
			G: 1,
			B: 1,
			A: 0x0f,
		}
		w.SetOffColor(off)
	}

	tdt.dt.Alignment = fyne.TextAlignCenter

	tdt.widget = widget.NewCard("", sidTheme.ToyDateTimeTitle,
		container.NewVBox(
			container.NewGridWithColumns(5, tdt.h1, tdt.h2, tdt.dot, tdt.m1, tdt.m2),
			tdt.dt),
	)

	// for init
	tdt.Run()

	_ = base.GlobalScheduler.AddJob("toy_datetime", "*/1 * * * * *", tdt)
	// todo assert(err == nil)
	return nil
}

func (tdt *toyDateTime) GetToyCard() *widget.Card {
	return tdt.widget
}

func (tdt *toyDateTime) Run() {
	if tdt.dot.String() == "" {
		tdt.dot.ParseMarkdown("# :")
	} else {

		tdt.dot.ParseMarkdown("")
	}

	now := time.Now()
	hr, min, _ := now.Clock()

	tdt.h1.Set(uint(hr / 10))
	tdt.h2.Set(uint(hr % 10))

	tdt.m1.Set(uint(min / 10))
	tdt.m2.Set(uint(min % 10))

	time.Sleep(time.Second)

	year, month, day := now.Date()

	var weekStr = [...]string{
		"Sunday",
		"Monday",
		"Tuesday",
		"Wednesday",
		"Thursday",
		"Friday",
		"Saturday",
	}

	tdt.dt.SetText(fmt.Sprintf("%d/%d/%d %s", year, month, day, weekStr[now.Weekday()]))
}
