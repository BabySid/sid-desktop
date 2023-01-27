package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	"sid-desktop/theme"
	"time"
)

type splashWin struct {
	loadProgress binding.Float
	loadStatus   binding.String
}

func (sw *splashWin) run() {
	drv := fyne.CurrentApp().Driver()
	if drv, ok := drv.(desktop.Driver); ok {
		w := drv.CreateSplashWindow()

		logo := canvas.NewImageFromResource(theme.ResourceSidLogo)
		logo.FillMode = canvas.ImageFillContain

		sw.loadProgress = binding.NewFloat()
		progressBar := widget.NewProgressBarWithData(sw.loadProgress)
		progressBar.Max = 100.0
		progressBar.Min = 0.0

		sw.loadStatus = binding.NewString()
		progressLabel := widget.NewLabelWithData(sw.loadStatus)

		w.SetContent(container.NewBorder(nil, container.NewVBox(progressLabel, progressBar), nil, nil, logo))
		w.Resize(fyne.NewSize(600, 500))
		w.Show()

		go func() {
			begin := 0.0
			for {
				sw.loadStatus.Set(gobase.FormatTimeStamp(time.Now().Unix()))
				sw.loadProgress.Set(begin)
				begin += 10
				if begin >= 100 {
					break
				}
				time.Sleep(time.Second * 3)
			}
			w.Hide()
		}()
	}
}
