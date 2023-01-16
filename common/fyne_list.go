package common

import (
	fyneTheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"sync"
)

func GetItemsHeightInList(list *widget.List, n int) float32 {
	var height float32
	if v, ok := listItemHeightCache.Load(list); ok {
		height = v.(float32)
	} else {
		height = list.CreateItem().MinSize().Height
		listItemHeightCache.Store(list, height)
	}

	listHeight := float32(n)*(height+2*fyneTheme.Padding()+fyneTheme.SeparatorThicknessSize()) + 2*fyneTheme.Padding()
	return listHeight
}

var (
	listItemHeightCache sync.Map
)
