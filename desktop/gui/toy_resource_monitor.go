package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	sidTheme "sid-desktop/desktop/theme"
)

var _ toyInterface = (*toyResourceMonitor)(nil)

type toyResourceMonitor struct {
	toyAdapter
	cpuIndicator *widget.ProgressBar
	memIndicator *widget.ProgressBar

	upTime *widget.Label
}

func (trm *toyResourceMonitor) Init() error {
	trm.cpuIndicator = widget.NewProgressBar()
	trm.cpuIndicator.Max = 100.0
	trm.cpuIndicator.Min = 0.0

	trm.memIndicator = widget.NewProgressBar()
	trm.memIndicator.Max = 100.0
	trm.memIndicator.Min = 0.0

	trm.upTime = widget.NewLabel("")
	trm.upTime.Alignment = fyne.TextAlignCenter

	trm.widget = widget.NewCard("", sidTheme.ToyResourceMonitorTitle,
		container.NewVBox(
			widget.NewForm(
				widget.NewFormItem(sidTheme.ToyResourceMonitorItem1, trm.cpuIndicator),
				widget.NewFormItem(sidTheme.ToyResourceMonitorItem2, trm.memIndicator),
			),
			trm.upTime,
		),
	)

	trm.widget.Resize(fyne.NewSize(ToyWidth, 150))

	trm.Run()

	_ = gobase.GlobalScheduler.AddJob("toy_resource_monitor", "*/1 * * * * *", trm)

	return nil
}

func (trm *toyResourceMonitor) Run() {
	trm.cpuIndicator.SetValue(gobase.GetCPUUsage())
	trm.memIndicator.SetValue(gobase.GetMEMUsage())

	uptime := gobase.GetUpTime()

	hr, min, sec := 0, 0, 0
	if uptime > 3600 {
		hr = int(uptime / 3600)
		mod := uptime % 3600
		if mod > 60 {
			min = int(mod / 60)
			sec = int(mod % 60)
		} else {
			min = 0
			sec = int(mod)
		}
	} else {
		hr = 0
		if uptime > 60 {
			min = int(uptime / 60)
			sec = int(uptime % 60)
		} else {
			min = 0
			sec = int(uptime)
		}
	}
	trm.upTime.SetText(fmt.Sprintf(sidTheme.ToyResourceMonitorUpTimeFormat, hr, min, sec))
}
