package common

import "fyne.io/fyne/v2/data/binding"

func CopyBindingStringList(dst binding.StringList, src binding.StringList) {
	srcList, _ := src.Get()
	temp := make([]string, len(srcList))
	copy(temp, srcList)
	_ = dst.Set(temp)
}

func Find(list binding.UntypedList, item binding.DataItem) int {
	length := list.Length()

	for i := 0; i < length; i++ {
		if data, _ := list.GetItem(i); data == item {
			return i
		}
	}

	return -1
}
