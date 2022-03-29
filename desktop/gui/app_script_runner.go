package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	"github.com/go-cmd/cmd"
	"github.com/satori/go.uuid"
	"google.golang.org/protobuf/encoding/protojson"
	"io/ioutil"
	"os"
	"path/filepath"
	"sid-desktop/desktop/common"
	"sid-desktop/desktop/storage"
	sidTheme "sid-desktop/desktop/theme"
	"sid-desktop/proto"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var _ appInterface = (*appScriptRunner)(nil)

type appScriptRunner struct {
	newScriptBtn   *widget.Button
	saveScriptBtn  *widget.Button
	stopScriptBtn  *widget.Button
	runScriptBtn   *widget.Button
	scriptLineNo   *widget.Label
	scriptBody     *widget.Entry
	logLabel       *widget.Label
	scriptLog      *widget.Entry
	scriptStatus   sync.Map // id -> *scriptStatus
	curScriptFile  *common.ScriptFile
	scriptBinding  binding.UntypedList
	scriptFiles    *widget.List
	scriptName     *widget.Entry
	runningScripts int32
	tabItem        *container.TabItem
}

type scriptStatus struct {
	Dirty bool
	Cmd   *cmd.Cmd
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
	asr.stopScriptBtn = widget.NewButtonWithIcon(sidTheme.AppScriptRunnerStopScript, sidTheme.ResourceStopIcon, asr.killScriptFile)

	asr.scriptLineNo = widget.NewLabel("1")
	asr.scriptBody = widget.NewMultiLineEntry()
	asr.scriptBody.OnChanged = func(s string) {
		if asr.curScriptFile == nil {
			return
		}
		if s == asr.curScriptFile.Cont {
			return
		}

		ln := strings.Count(s, "\n")
		nu := ""
		for i := 1; i <= ln+1; i++ {
			nu += fmt.Sprintf("%d\n", i)
		}

		asr.scriptLineNo.SetText(nu)

		asr.setCurScriptDirty(asr.curScriptFile.ID, true)
	}

	asr.logLabel = widget.NewLabel(sidTheme.AppScriptRunnerRunLog)
	asr.scriptLog = widget.NewMultiLineEntry()

	asr.scriptName = widget.NewEntry()
	asr.scriptName.Validator = validation.NewRegexp(`\S+`, sidTheme.AppScriptRunnerScriptNameValidateMsg)
	asr.scriptName.SetPlaceHolder(sidTheme.AppScriptRunnerCurScriptName)
	asr.scriptName.OnChanged = func(s string) {
		if asr.curScriptFile == nil {
			return
		}
		if s == asr.curScriptFile.Name {
			return
		}
		if !strings.HasSuffix(s, ".lua") {
			// Support .lua only now
			asr.scriptName.SetText(s + ".lua")
		}

		asr.setCurScriptDirty(asr.curScriptFile.ID, true)
	}

	asr.scriptBinding = binding.NewUntypedList()
	asr.createScriptList()
	asr.scriptFiles.OnSelected = func(id widget.ListItemID) {
		sf, _ := asr.scriptBinding.GetValue(id)
		file, _ := sf.(common.ScriptFile)

		if asr.curScriptFile != nil && asr.curScriptFile.ID != file.ID {
			status := asr.getScriptStatus(asr.curScriptFile.ID)
			if status != nil && status.Dirty {
				asr.saveScriptFile()
				asr.setCurScriptDirty(asr.curScriptFile.ID, false)
			}
		}

		asr.curScriptFile = &file
		asr.scriptName.SetText(asr.curScriptFile.Name)
		asr.scriptBody.SetText(asr.curScriptFile.Cont)
	}

	asr.tabItem = container.NewTabItemWithIcon(sidTheme.AppScriptRunnerName, sidTheme.ResourceScriptRunnerIcon, nil)

	scriptPanel := container.NewVSplit(
		container.NewBorder(
			container.NewGridWithColumns(2,
				asr.scriptName,
				container.NewHBox(layout.NewSpacer(), asr.saveScriptBtn, asr.stopScriptBtn, asr.runScriptBtn)),
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
	v := atomic.LoadInt32(&asr.runningScripts)
	if v != 0 {
		dialog.ShowInformation(sidTheme.CannotCloseTitle, sidTheme.AppScriptRunnerCannotCloseMsg, globalWin.win)
		return false
	}

	return true
}

func (asr *appScriptRunner) createScriptList() {
	// Script List
	asr.scriptFiles = widget.NewListWithData(
		asr.scriptBinding,
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(sidTheme.ResourceLuaIcon),
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				layout.NewSpacer(),
				widget.NewButtonWithIcon(sidTheme.AppScriptRunnerDelScript, sidTheme.ResourceRmIcon, nil))
		},
		func(data binding.DataItem, item fyne.CanvasObject) {
			o, _ := data.(binding.Untyped).Get()
			file := o.(common.ScriptFile)

			name := file.Name
			status := asr.getScriptStatus(file.ID)
			if status != nil {
				if status.Cmd != nil {
					name += "(R)"
				}
				if status.Dirty {
					name += "*"
				}
			}
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(name)

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
	txt := strings.Trim(asr.scriptName.Text, "\t ")
	if txt == "" {
		dialog.ShowInformation(sidTheme.AppScriptRunnerName, sidTheme.AppScriptRunnerScriptNameValidateMsg, globalWin.win)
		return
	}
	if asr.curScriptFile == nil {
		asr.curScriptFile = &common.ScriptFile{
			Name:       txt,
			Cont:       asr.scriptBody.Text,
			CreateTime: time.Now().Unix(),
			AccessTime: time.Now().Unix(),
		}
		err := storage.GetAppScriptRunnerDB().AddScriptFile(*asr.curScriptFile)
		if err != nil {
			printErr(fmt.Errorf(sidTheme.ProcessScriptRunnerFailedFormat, err))
			return
		}
	} else {
		asr.curScriptFile.Name = txt
		asr.curScriptFile.Cont = asr.scriptBody.Text
		asr.curScriptFile.AccessTime = time.Now().Unix()
		err := storage.GetAppScriptRunnerDB().UpdateScriptFile(*asr.curScriptFile)
		if err != nil {
			printErr(fmt.Errorf(sidTheme.ProcessScriptRunnerFailedFormat, err))
			return
		}
	}

	asr.setCurScriptDirty(asr.curScriptFile.ID, false)
	asr.reloadScriptFiles()
}

func (asr *appScriptRunner) killScriptFile() {
	if asr.curScriptFile == nil {
		return
	}
	status := asr.getScriptStatus(asr.curScriptFile.ID)
	if status == nil || status.Cmd == nil {
		return
	}

	_ = status.Cmd.Stop()

	asr.setCurScriptRunningHandle(asr.curScriptFile.ID, nil)
}

func (asr *appScriptRunner) runScriptFile() {
	if asr.curScriptFile != nil {
		status := asr.getScriptStatus(asr.curScriptFile.ID)
		if status != nil && status.Dirty {
			asr.saveScriptFile()
		}
	}

	if asr.curScriptFile == nil {
		return
	}

	scriptFile := *asr.curScriptFile

	tmpFilePath := filepath.Join(os.TempDir(), uuid.NewV4().String())
	msg := &proto.ScriptRunner{
		Id:      scriptFile.ID,
		Title:   scriptFile.Name,
		Content: scriptFile.Cont,
	}
	script, err := protojson.Marshal(msg)
	if err != nil {
		printErr(fmt.Errorf(sidTheme.ProcessScriptRunnerFailedFormat, err))
		return
	}
	err = ioutil.WriteFile(tmpFilePath, []byte(script), 0600)
	if err != nil {
		printErr(fmt.Errorf(sidTheme.ProcessScriptRunnerFailedFormat, err))
		return
	}

	c := cmd.NewCmdOptions(cmd.Options{
		Buffered:   false,
		Streaming:  true,
		BeforeExec: nil,
	}, common.GetLuaRunner(), "--script", tmpFilePath)

	asr.setCurScriptRunningHandle(scriptFile.ID, c)

	s := c.Start()

	go func() {
		atomic.AddInt32(&asr.runningScripts, 1)
		defer func() {
			atomic.AddInt32(&asr.runningScripts, -1)
			asr.setCurScriptRunningHandle(scriptFile.ID, nil)
			asr.reloadScriptFiles()
		}()

		for {
			select {
			case v := <-c.Stdout:
				asr.appendScriptLog(v)
			case <-c.Done():
				return
			}
		}
	}()

	status := <-s
	asr.appendScriptLog(fmt.Sprintf("script [%s] exit with %d\n", scriptFile.Name, status.Exit))
}

func (asr *appScriptRunner) appendScriptLog(newLog string) {
	txt := asr.scriptLog.Text + "\n"
	txt += newLog
	asr.scriptLog.CursorRow = strings.Count(txt, "\n")
	asr.scriptLog.SetText(txt)
}

func (asr *appScriptRunner) setCurScriptDirty(id int32, dirty bool) {
	status := asr.getScriptStatus(id)
	if status == nil {
		status = &scriptStatus{
			Dirty: dirty,
			Cmd:   nil,
		}
	}

	asr.scriptStatus.Store(id, status)
}

func (asr *appScriptRunner) setCurScriptRunningHandle(id int32, c *cmd.Cmd) {
	status := asr.getScriptStatus(id)
	if status == nil {
		status = &scriptStatus{
			Dirty: false,
			Cmd:   c,
		}
	}

	asr.scriptStatus.Store(id, status)
}

func (asr *appScriptRunner) getScriptStatus(id int32) *scriptStatus {
	if v, ok := asr.scriptStatus.Load(id); ok {
		return v.(*scriptStatus)
	}

	return nil
}
