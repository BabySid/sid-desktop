package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var _ sodorInterface = (*sodorAlertPlugin)(nil)

type sodorAlertPlugin struct {
	sodorAdapter
}

func (s *sodorAlertPlugin) CreateView() fyne.CanvasObject {
	if s.content != nil {
		return s.content
	}

	s.content = container.NewBorder(nil,
		nil, nil, nil, widget.NewLabel("alert plugin"))
	return s.content
}
