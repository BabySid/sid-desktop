package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	"path/filepath"
	"sid-desktop/desktop/common"
	"sid-desktop/desktop/storage"
	sidTheme "sid-desktop/desktop/theme"
	"strings"
	"time"
)

var _ appInterface = (*appScriptRunner)(nil)

type appScriptRunner struct {
	newScriptBtn  *widget.Button
	saveScriptBtn *widget.Button
	runScriptBtn  *widget.Button
	scriptLineNo  *widget.Label
	scriptBody    *widget.Entry
	logLabel      *widget.Label
	scriptLog     *widget.Entry
	curScriptFile *common.ScriptFile
	scriptBinding binding.UntypedList
	scriptFiles   *widget.List
	scriptName    *widget.Entry
	tabItem       *container.TabItem
}

func (asr *appScriptRunner) LazyInit() error {
	err := storage.GetAppScriptRunnerDB().Open(globalWin.app.Storage().RootURI().Path())
	if err != nil {
		return err
	}
	gobase.RegisterAtExit(storage.GetAppScriptRunnerDB().Close)

	asr.newScriptBtn = widget.NewButtonWithIcon(sidTheme.AppScriptRunnerNewScript, sidTheme.ResourceAddIcon, asr.newScriptFile)
	asr.saveScriptBtn = widget.NewButtonWithIcon(sidTheme.AppScriptRunnerSaveScript, sidTheme.ResourceSaveIcon, asr.saveScriptFile)
	asr.runScriptBtn = widget.NewButtonWithIcon(sidTheme.AppScriptRunnerRunScript, sidTheme.ResourceRunIcon, asr.runScriptFile)

	asr.scriptLineNo = widget.NewLabel("1")
	asr.scriptBody = widget.NewMultiLineEntry()
	asr.scriptBody.OnChanged = func(s string) {
		ln := strings.Count(s, "\n")
		nu := ""
		for i := 1; i <= ln+1; i++ {
			nu += fmt.Sprintf("%d\n", i)
		}

		asr.scriptLineNo.SetText(nu)

		if asr.curScriptFile != nil {
			asr.curScriptFile.Dirty = true
		}
	}

	asr.logLabel = widget.NewLabel(sidTheme.AppScriptRunnerRunLog)
	asr.scriptLog = widget.NewMultiLineEntry()

	asr.scriptName = widget.NewEntry()
	asr.scriptName.Validator = validation.NewRegexp(`\S+`, sidTheme.AppScriptRunnerCurScriptName+" must not be empty")
	asr.scriptName.SetPlaceHolder(sidTheme.AppScriptRunnerCurScriptName)
	asr.scriptName.OnChanged = func(s string) {
		if !strings.HasSuffix(s, ".lua") {
			// Support .lua only nnow
			asr.scriptName.SetText(s + ".lua")
		}

		if asr.curScriptFile != nil {
			asr.curScriptFile.Dirty = true
		}
	}

	asr.scriptBinding = binding.NewUntypedList()
	asr.createScriptList()
	asr.scriptFiles.OnSelected = func(id widget.ListItemID) {
		if asr.curScriptFile != nil && asr.curScriptFile.Dirty {
			asr.saveScriptFile()
		}

		sf, _ := asr.scriptBinding.GetValue(id)
		file, _ := sf.(common.ScriptFile)
		asr.curScriptFile = &file
		asr.scriptName.SetText(asr.curScriptFile.Name)
		asr.scriptBody.SetText(asr.curScriptFile.Cont)
	}

	asr.tabItem = container.NewTabItemWithIcon(sidTheme.AppScriptRunnerName, sidTheme.ResourceScriptRunnerIcon, nil)

	scriptPanel := container.NewVSplit(
		container.NewBorder(
			container.NewGridWithColumns(2,
				asr.scriptName,
				container.NewHBox(layout.NewSpacer(), asr.saveScriptBtn, asr.runScriptBtn)),
			nil, asr.scriptLineNo, nil,
			asr.scriptBody),
		container.NewBorder(container.NewHBox(asr.logLabel, layout.NewSpacer()), nil, nil, nil,
			asr.scriptLog),
	)
	scriptPanel.SetOffset(0.8)

	contPanel := container.NewHSplit(container.NewBorder(container.NewHBox(layout.NewSpacer(), asr.newScriptBtn),
		nil, nil, nil,
		asr.scriptFiles), scriptPanel)
	contPanel.SetOffset(0.3)

	asr.tabItem.Content = contPanel

	go asr.initDB()

	return nil
}

func (asr *appScriptRunner) GetTabItem() *container.TabItem {
	return asr.tabItem
}

func (asr *appScriptRunner) GetAppName() string {
	return sidTheme.AppScriptRunnerName
}

func (asr *appScriptRunner) OpenDefault() bool {
	return false
}

func (asr *appScriptRunner) OnClose() bool {
	// todo run status
	return true
}

func (asr *appScriptRunner) createScriptList() {
	// Script List
	asr.scriptFiles = widget.NewListWithData(
		asr.scriptBinding,
		func() fyne.CanvasObject {
			// todo lua icon
			return container.NewHBox(
				widget.NewIcon(sidTheme.ResourceLuaIcon),
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				layout.NewSpacer(),
				widget.NewButtonWithIcon(sidTheme.AppScriptRunnerDelScript, sidTheme.ResourceRmIcon, nil))
		},
		func(data binding.DataItem, item fyne.CanvasObject) {
			o, _ := data.(binding.Untyped).Get()
			file := o.(common.ScriptFile)

			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(file.Name)

			item.(*fyne.Container).Objects[3].(*widget.Button).SetText(sidTheme.AppScriptRunnerDelScript)
			item.(*fyne.Container).Objects[3].(*widget.Button).OnTapped = func() {
				err := storage.GetAppScriptRunnerDB().DelScriptFile(file)
				if err != nil {
					printErr(fmt.Errorf(sidTheme.ProcessScriptRunnerFailedFormat, err))
					return
				}

				if asr.curScriptFile != nil && asr.curScriptFile.ID == file.ID {
					asr.curScriptFile = nil
				}

				asr.reloadScriptFiles()
			}
		},
	)
}

func (asr *appScriptRunner) initDB() {
	need, err := storage.GetAppScriptRunnerDB().NeedInit()
	if err != nil {
		printErr(fmt.Errorf(sidTheme.ProcessScriptRunnerFailedFormat, err))
		return
	}

	if need {
		err = storage.GetAppScriptRunnerDB().Init()
		if err != nil {
			printErr(fmt.Errorf(sidTheme.ProcessScriptRunnerFailedFormat, err))
			return
		}
	} else {
		asr.reloadScriptFiles()
	}
}

func (asr *appScriptRunner) reloadScriptFiles() {
	asr.loadScriptFilesFromDB()
}

func (asr *appScriptRunner) loadScriptFilesFromDB() {
	files, err := storage.GetAppScriptRunnerDB().LoadScriptFiles()
	if err != nil {
		printErr(fmt.Errorf(sidTheme.ProcessScriptRunnerFailedFormat, err))
		return
	}

	_ = asr.scriptBinding.Set(files.AsInterfaceArray())
}

func (asr *appScriptRunner) newScriptFile() {
	docID := 0
	files, _ := asr.scriptBinding.Get()
	for _, file := range files {
		if filepath.Base(file.(common.ScriptFile).Name) == sidTheme.AppScriptRunnerNewScriptName {
			docID++
			fmt.Println(docID)
		}
	}

	if docID == 0 {
		asr.scriptName.SetText(sidTheme.AppScriptRunnerNewScriptName)
	} else {
		asr.scriptName.SetText(sidTheme.AppScriptRunnerNewScriptName + fmt.Sprintf("%d", docID))
	}

	asr.curScriptFile = nil
	asr.scriptBody.SetText("")
	asr.scriptLog.SetText("")
}

func (asr *appScriptRunner) saveScriptFile() {
	if asr.curScriptFile == nil {
		asr.curScriptFile = &common.ScriptFile{
			Name:       asr.scriptName.Text,
			Cont:       asr.scriptBody.Text,
			CreateTime: time.Now().Unix(),
			AccessTime: time.Now().Unix(),
		}
		err := storage.GetAppScriptRunnerDB().AddScriptFile(*asr.curScriptFile)
		if err != nil {
			printErr(fmt.Errorf(sidTheme.ProcessScriptRunnerFailedFormat, err))
		}
	} else {
		asr.curScriptFile.Name = asr.scriptName.Text
		asr.curScriptFile.Cont = asr.scriptBody.Text
		asr.curScriptFile.AccessTime = time.Now().Unix()
		err := storage.GetAppScriptRunnerDB().UpdateScriptFile(*asr.curScriptFile)
		if err != nil {
			printErr(fmt.Errorf(sidTheme.ProcessScriptRunnerFailedFormat, err))
		}
	}

	asr.curScriptFile.Dirty = false
	asr.reloadScriptFiles()
}

func (asr *appScriptRunner) runScriptFile() {
	if asr.curScriptFile != nil && asr.curScriptFile.Dirty {
		asr.saveScriptFile()
	}

	// todo run
	fmt.Println(asr.curScriptFile.Name)
}

func (asr *appScriptRunner) setCurScriptDirty(dirty bool) {
	if asr.curScriptFile != nil {
		asr.curScriptFile.Dirty = dirty
	}
}
