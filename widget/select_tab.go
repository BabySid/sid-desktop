package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type SelectTab struct {
	sel     *widget.Select
	Content fyne.CanvasObject

	OnSelected func(*SelectItem)
}

func NewSelectTab(items ...*SelectItem) *SelectTab {
	options := make([]string, len(items))
	itemContents := make([]fyne.CanvasObject, len(items))
	for i, item := range items {
		options[i] = item.Text
		itemContents[i] = item.Content
	}

	st := &SelectTab{}

	st.sel = widget.NewSelect(options, func(s string) {
		for _, item := range items {
			if s == item.Text {
				item.Content.Show()

				if st.OnSelected != nil {
					st.OnSelected(item)
				}
			} else {
				item.Content.Hide()
			}
		}
	})
	st.sel.SetSelectedIndex(0)

	content := container.NewBorder(
		container.NewHBox(st.sel, layout.NewSpacer()), nil, nil, nil, itemContents...)
	st.Content = content

	return st
}

func (st *SelectTab) SetSelected(s string) {
	st.sel.SetSelected(s)
}

type SelectItem struct {
	Text    string
	Content fyne.CanvasObject
}

func NewSelectItem(name string, cont fyne.CanvasObject) *SelectItem {
	return &SelectItem{
		Text:    name,
		Content: cont,
	}
}
