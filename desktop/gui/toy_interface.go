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
