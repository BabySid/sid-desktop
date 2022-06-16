package gui

import (
	"fyne.io/fyne/v2/widget"
)

type toyInterface interface {
	Init() error
	GetToyCard() *widget.Card
}

var (
	toyRegister = []toyInterface{
		&toyResourceMonitor{},
		&toyDateTime{},
		&toyHotSearch{},
	}
)

const (
	ToyWidth = 250
)

var _ toyInterface = (*toyAdapter)(nil)

type toyAdapter struct {
	widget *widget.Card
}

func (t toyAdapter) Init() error {
	panic("implement Init")
}

func (t toyAdapter) GetToyCard() *widget.Card {
	return t.widget
}
