package common

import (
	fyneTheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"sync"
)

func GetItemsHeightInCheck(n int) float32 {
	var height float32
	if v, ok := checkItemHeightCache.Load(checkGroupKey); ok {
		height = v.(float32)
	} else {
		chk := widget.NewCheckGroup([]string{"Test"}, nil)
		height = chk.MinSize().Height
		checkItemHeightCache.Store(checkGroupKey, height)
	}

	cgHeight := float32(n)*(height+2*fyneTheme.Padding()+fyneTheme.SeparatorThicknessSize()) + 2*fyneTheme.Padding()
	return cgHeight
}

var (
	checkItemHeightCache sync.Map
	checkGroupKey        = "checkGroup"
)
