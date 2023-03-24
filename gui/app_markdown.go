package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	"sid-desktop/common"
	"sid-desktop/storage"
	"sid-desktop/theme"
)

var _ appInterface = (*appMarkDown)(nil)

type appMarkDown struct {
	appAdapter

	fileList      *widget.List
	fileBinding   binding.UntypedList
	searchEntry   *widget.Entry
	newFile       *widget.Button
	saveFile      *widget.Button
	fileNameEntry *widget.Entry
	editEntry     *widget.Entry
	previewEntry  *widget.RichText

	fileCache *common.MarkDownFileList

	curFile *common.MarkDownFile
}

func (amd *appMarkDown) LazyInit() error {
	err := storage.GetAppMarkDownDB().Open(globalWin.app.Storage().RootURI().Path())
	if err != nil {
		return err
	}
	gobase.RegisterAtExit(storage.GetAppMarkDownDB().Close)

	amd.searchEntry = widget.NewEntry()
	amd.searchEntry.SetPlaceHolder("Search Title and Content")
	amd.searchEntry.OnChanged = amd.searchMarkDown

	amd.saveFile = widget.NewButtonWithIcon("Save", theme.ResourceSaveIcon, amd.saveMarkDownFile)
	amd.newFile = widget.NewButtonWithIcon("New", theme.ResourceAddIcon, amd.addFile)

	amd.fileBinding = binding.NewUntypedList()
	amd.createFileList()

	amd.fileNameEntry = widget.NewEntry()
	amd.editEntry = widget.NewMultiLineEntry()
	amd.editEntry.PlaceHolder = `# 一级标题
## 二级标题
### 三级标题

### 无序列表

- 第一项
- 第二项
- 第三项

### 有序列表

1. 第一项
2. 第二项
3. 第三项

> 这是一段引用的内容。

*斜体*
**粗体**
***粗斜体***
` + "```" + `
func main() {
	fmt.Println("Hello World!")
}
` + "```" + `
Markdown中使用[]和()符号表示链接，例如：
[百度](https://www.baidu.com/)
	`
	amd.editEntry.OnChanged = func(s string) {
		amd.previewEntry.ParseMarkdown(s)
	}
	amd.previewEntry = widget.NewRichTextFromMarkdown("")

	amd.tabItem = container.NewTabItemWithIcon("MarkDown", theme.ResourceMarkDownIcon, nil)
	split := container.NewHSplit(amd.fileList, container.NewHSplit(
		container.NewBorder(amd.fileNameEntry, nil, nil, nil, amd.editEntry), amd.previewEntry))
	split.SetOffset(0.2)

	amd.tabItem.Content = container.NewBorder(
		container.NewGridWithColumns(3, amd.searchEntry, layout.NewSpacer(),
			container.NewHBox(layout.NewSpacer(), amd.newFile, amd.saveFile)),
		nil, nil, nil,
		split,
	)

	//go amd.initDB()

	return nil
}

func (amd *appMarkDown) GetAppName() string {
	return "MarkDown"
}

func (amd *appMarkDown) OnClose() bool {
	return true
}

func (amd *appMarkDown) ShortCut() fyne.Shortcut {
	return &desktop.CustomShortcut{KeyName: fyne.Key4, Modifier: fyne.KeyModifierAlt}
}

func (amd *appMarkDown) Icon() fyne.Resource {
	return theme.ResourceMarkDownIcon
}

func (amd *appMarkDown) searchMarkDown(s string) {

}

func (amd *appMarkDown) addFile() {
	amd.fileNameEntry.SetText("")
	amd.editEntry.SetText("")
}

func (amd *appMarkDown) createFileList() {
	amd.fileList = widget.NewListWithData(
		amd.fileBinding,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		}, func(data binding.DataItem, item fyne.CanvasObject) {
			o, _ := data.(binding.Untyped).Get()
			md := o.(string)
			item.(*widget.Label).SetText(md)
		})

	amd.fileBinding.Set([]interface{}{"file1", "file2", "file3"})
}

func (amd *appMarkDown) saveMarkDownFile() {
	file := common.MarkDownFile{
		Name: amd.fileNameEntry.Text,
		Cont: amd.editEntry.Text,
	}
	err := storage.GetAppMarkDownDB().UpdateMarkDownFile(&file)
	if err != nil {
		printErr(fmt.Errorf("markdown failed %v", err))
	}
}

func (amd *appMarkDown) initDB() {
	need, err := storage.GetAppMarkDownDB().NeedInit()
	if err != nil {
		printErr(fmt.Errorf("markdown failed %v", err))
		return
	}

	if need {
		err = storage.GetAppMarkDownDB().Init()
		if err != nil {
			printErr(fmt.Errorf("markdown failed %v", err))
			return
		}
	} else {
		amd.reloadMarkDownFiles()
	}
}

func (amd *appMarkDown) reloadMarkDownFiles() {
	var err error
	amd.fileCache, err = storage.GetAppMarkDownDB().LoadMarkDownFiles()
	if err != nil {
		printErr(fmt.Errorf("markdown failed %v", err))
	}
	if amd.fileCache != nil {
		_ = amd.fileBinding.Set(amd.fileCache.AsInterfaceArray())
	}
}
