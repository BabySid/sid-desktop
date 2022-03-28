package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	"sid-desktop/desktop/common"
	"sid-desktop/desktop/storage"
	sidTheme "sid-desktop/desktop/theme"
	"strings"
)

var _ appInterface = (*appScriptRunner)(nil)

type appScriptRunner struct {
	newScript     *widget.Button
	runScript     *widget.Button
	scriptLineNo  *widget.Label
	scriptText    *widget.Entry
	logLabel      *widget.Label
	scriptLog     *widget.Entry
	scriptBinding binding.UntypedList
	scriptFiles   *widget.List
	curScript     *widget.Entry
	tabItem       *container.TabItem
}

func (asr *appScriptRunner) LazyInit() error {
	err := storage.GetAppScriptRunnerDB().Open(globalWin.app.Storage().RootURI().Path())
	if err != nil {
		return err
	}
	gobase.RegisterAtExit(storage.GetAppScriptRunnerDB().Close)

	asr.newScript = widget.NewButton(sidTheme.AppScriptRunnerNewScript, asr.newScriptFile)
	asr.runScript = widget.NewButton(sidTheme.AppScriptRunnerRunScript, nil)

	asr.scriptLineNo = widget.NewLabel("1")
	asr.scriptText = widget.NewMultiLineEntry()
	asr.scriptText.OnChanged = func(s string) {
		ln := strings.Count(s, "\n")
		nu := ""
		for i := 1; i <= ln+1; i++ {
			nu += fmt.Sprintf("%d\n", i)
		}

		asr.scriptLineNo.SetText(nu)
	}

	asr.logLabel = widget.NewLabel(sidTheme.AppScriptRunnerRunLog)
	asr.scriptLog = widget.NewMultiLineEntry()

	asr.curScript = widget.NewEntry()
	asr.curScript.SetText(fmt.Sprintf(sidTheme.AppScriptRunnerCurScriptFormat, ""))
	asr.curScript.OnChanged = func(s string) {
		if !strings.HasSuffix(s, ".lua") {
			// now we support .lua only
			asr.curScript.SetText(s + ".lua")
		}
	}

	asr.scriptBinding = binding.NewUntypedList()
	asr.createScriptList()
	asr.scriptFiles.OnSelected = func(id widget.ListItemID) {
		sf, _ := asr.scriptBinding.GetValue(id)
		file := sf.(common.ScriptFile)
		asr.curScript.SetText(file.Name)

		asr.scriptText.SetText(file.Cont)
	}

	asr.tabItem = container.NewTabItemWithIcon(sidTheme.AppScriptRunnerName, sidTheme.ResourceScriptRunnerIcon, nil)

	scriptPanel := container.NewVSplit(
		container.NewBorder(
			container.NewGridWithColumns(2,
				asr.curScript,
				container.NewHBox(layout.NewSpacer(), asr.runScript)),
			nil, asr.scriptLineNo, nil,
			asr.scriptText),
		container.NewBorder(container.NewHBox(asr.logLabel, layout.NewSpacer()), nil, nil, nil,
			asr.scriptLog),
	)
	scriptPanel.SetOffset(0.8)

	contPanel := container.NewHSplit(container.NewBorder(container.NewHBox(layout.NewSpacer(), asr.newScript),
		nil, nil, nil,
		asr.scriptFiles), scriptPanel)
	contPanel.SetOffset(0.2)

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
	return true
}

func (asr *appScriptRunner) createScriptList() {
	// Script List
	asr.scriptFiles = widget.NewListWithData(
		asr.scriptBinding,
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				widget.NewButton(sidTheme.AppScriptRunnerDelScript, nil))
		},
		func(data binding.DataItem, item fyne.CanvasObject) {
			o, _ := data.(binding.Untyped).Get()
			file := o.(common.ScriptFile)

			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(file.Name)

			item.(*fyne.Container).Objects[1].(*widget.Button).SetText(sidTheme.AppScriptRunnerDelScript)
			item.(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {

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
	asr.curScript.SetText(sidTheme.AppScriptRunnerNewScriptName)
	asr.scriptText.SetText("")
	asr.scriptLog.SetText("")
}
