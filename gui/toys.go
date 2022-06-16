package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	sidContainer "sid-desktop/container"
)

type toys struct {
	widget *fyne.Container
}

func newToys() *toys {
	var t toys

	cards := make([]fyne.CanvasObject, len(toyRegister))
	for i, toy := range toyRegister {
		toy := toy
		err := toy.Init()
		if err != nil {
			panic(err)
		}

		cards[i] = toy.GetToyCard()
	}

	t.widget = container.New(sidContainer.NewVCZBoxLayout(),
		cards...)
	//t.widget = container.NewVBox(container.NewGridWithRows(len(toyRegister), cards...),
	//	layout.NewSpacer())

	return &t
}
