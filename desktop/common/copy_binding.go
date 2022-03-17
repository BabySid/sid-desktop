package common

import "fyne.io/fyne/v2/data/binding"

func CopyBindingStringList(dst binding.StringList, src binding.StringList) {
	srcList, _ := src.Get()
	temp := make([]string, len(srcList))
	copy(temp, srcList)
	_ = dst.Set(temp)
}
