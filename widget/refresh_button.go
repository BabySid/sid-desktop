package widget

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	"sid-desktop/theme"
)

type RefreshButton struct {
	refreshBtn  *widget.Button
	autoRefresh *widget.Check

	Content fyne.CanvasObject

	routineName string
	running     bool
	OnRefresh   func()
}

func NewRefreshButton(spec string, tapped func()) *RefreshButton {
	btn := RefreshButton{}

	btn.OnRefresh = tapped
	btn.refreshBtn = widget.NewButtonWithIcon(theme.AppPageRefresh, theme.ResourceRefreshIcon, btn.OnRefresh)

	btn.routineName = fmt.Sprintf(refreshButtonRoutine, &btn)

	btn.autoRefresh = widget.NewCheck(theme.AppPageAutoRefresh, func(b bool) {
		if btn.running {
			gobase.GlobalScheduler.DelJob(btn.routineName)
		}

		if !b {
			btn.running = false
			return
		}

		btn.running = true
		gobase.GlobalScheduler.AddJob(btn.routineName, spec, &btn)
	})
	btn.autoRefresh.SetChecked(false)

	btn.Content = container.NewHBox(btn.refreshBtn, btn.autoRefresh)

	return &btn
}

const refreshButtonRoutine = "refresh_button_routine_%p"

func (rb *RefreshButton) Run() {
	if rb.OnRefresh != nil {
		rb.OnRefresh()
	}
}

func (rb *RefreshButton) Stop() {
	if rb.running {
		gobase.GlobalScheduler.DelJob(rb.routineName)
	}
	rb.running = false
}
