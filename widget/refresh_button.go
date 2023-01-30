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
	refreshBtn *widget.Button
	refreshSel *widget.Select

	Content fyne.CanvasObject

	routineName string
	running     bool
	OnRefresh   func()
}

func NewRefreshButton(icon fyne.Resource, tapped func()) *RefreshButton {
	btn := RefreshButton{}

	btn.OnRefresh = tapped
	if icon == nil {
		btn.refreshBtn = widget.NewButton(theme.AppPageRefresh, btn.OnRefresh)
	} else {
		btn.refreshBtn = widget.NewButtonWithIcon(theme.AppPageRefresh, icon, btn.OnRefresh)
	}

	btn.routineName = fmt.Sprintf(refreshButtonRoutine, &btn)

	btn.refreshSel = widget.NewSelect(btn.getCheckInterval(), func(s string) {
		if btn.running {
			gobase.GlobalScheduler.DelJob(btn.routineName)
		}

		if s == disableAuto {
			btn.running = false
			return
		}

		btn.running = true
		gobase.GlobalScheduler.AddJob(btn.routineName, checkIntervalToCronSpec[s], &btn)
	})
	btn.refreshSel.PlaceHolder = disableAuto
	btn.refreshSel.SetSelectedIndex(0)

	btn.Content = container.NewHBox(btn.refreshBtn, btn.refreshSel)

	return &btn
}

const disableAuto = "--NotAutoRefresh--"
const refreshButtonRoutine = "refresh_button_routine_%p"

var checkIntervalToCronSpec = map[string]string{
	"30s": "*/30 * * * * *",
	"1m":  "0 */1 * * * *",
	"5m":  "0 */5 * * * *",
}

func (rb *RefreshButton) getCheckInterval() []string {
	return []string{
		disableAuto,
		"30s",
		"1m",
		"5m",
	}
}

func (rb *RefreshButton) Run() {
	if rb.OnRefresh != nil {
		rb.OnRefresh()
	}
}
