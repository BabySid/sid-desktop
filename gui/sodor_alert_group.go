package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var _ sodorInterface = (*sodorAlertGroup)(nil)

type sodorAlertGroup struct {
	sodorAdapter
}

func (s *sodorAlertGroup) CreateView() fyne.CanvasObject {
	if s.content != nil {
		return s.content
	}

	s.content = container.NewBorder(nil,
		nil, nil, nil, widget.NewLabel("alert group"))
	return s.content
}
