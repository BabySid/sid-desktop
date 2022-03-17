package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"image/color"
	"time"
)

type splashWin struct {
}

func (sw *splashWin) run() {
	drv := fyne.CurrentApp().Driver()
	if drv, ok := drv.(desktop.Driver); ok {
		w := drv.CreateSplashWindow()

		line := canvas.NewLine(color.NRGBA{R: 0, G: 0, B: 180, A: 128})
		line.StrokeWidth = 5
		line.Position1 = fyne.NewPos(0, 0)
		line.Position2 = fyne.NewPos(50, 0)
		w.SetContent(container.NewWithoutLayout(line))
		w.Resize(fyne.NewSize(800, 3))
		w.Show()

		a2 := canvas.NewPositionAnimation(fyne.NewPos(0, 0), fyne.NewPos(800, 0), time.Second*3, func(p fyne.Position) {
			line.Move(p)
			line.Refresh()
		})
		a2.RepeatCount = fyne.AnimationRepeatForever
		a2.AutoReverse = true
		a2.Curve = fyne.AnimationLinear
		a2.Start()
	}
}
