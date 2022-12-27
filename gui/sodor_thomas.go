package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var _ sodorInterface = (*sodorThomas)(nil)

type sodorThomas struct {
	sodorAdapter
}

func (s *sodorThomas) CreateView() fyne.CanvasObject {
	if s.content != nil {
		return s.content
	}

	s.content = container.NewBorder(nil,
		nil, nil, nil, widget.NewLabel("thomas"))
	return s.content
}
