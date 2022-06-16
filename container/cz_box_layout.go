package container

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// Declare conformity with Layout interface
var _ fyne.Layout = (*CZBoxLayout)(nil)

// CZBoxLayout is a grid layout that support custom size of object
// Now only support vertical mode
type CZBoxLayout struct {
	vertical bool
}

func NewVCZBoxLayout() fyne.Layout {
	return &CZBoxLayout{vertical: true}
}

func (c *CZBoxLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	pos := fyne.NewPos(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		child.Move(pos)
		size := child.Size()
		if c.vertical {
			pos = pos.Add(fyne.NewPos(0, size.Height+theme.Padding()))
		}
	}
}

func (c *CZBoxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	addPadding := false
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if c.vertical {
			size := child.Size()
			minSize.Width = fyne.Max(size.Width, minSize.Width)
			minSize.Height += size.Height
			if addPadding {
				minSize.Height += theme.Padding()
			}
		}

		addPadding = true
	}

	return minSize
}
