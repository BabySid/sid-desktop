package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"log"
	"sid-desktop/base"
	"sid-desktop/desktop/common"
	sidTheme "sid-desktop/desktop/theme"
	"time"
)

type MainWin struct {
	app fyne.App
	win fyne.Window

	mm   *mainMenu
	tb   *toolBar
	toys *toys
	sb   *statusBar
	at   *appContainer
}

func init() {
	// set env to support chinese
	//_ = os.Setenv("FYNE_FONT", "./resource/Microsoft-YaHei.ttf")
	//_ = os.Setenv("FYNE_FONT_MONOSPACE", "./resource/Microsoft-YaHei.ttf")
}

//func setAPPID() {
//	// SetCurrentProcessExplicitAppUserModelID
//	shell32 := syscall.NewLazyDLL("shell32.dll")
//	fun := shell32.NewProc("SetCurrentProcessExplicitAppUserModelID")
//	appid := "Sid desktop"
//	p, err := syscall.UTF16PtrFromString(appid)
//	if err != nil {
//		panic(err)
//	}
//	in := uintptr(unsafe.Pointer(p))
//	r1, r2, e := fun.Call(in)
//	fmt.Println(r1, r2, e)
//}

var (
	globalWin       *MainWin
	globalConfig    *common.Config
	globalLogWriter *common.LogWriter
)

func NewMainWin() *MainWin {
	base.NewScheduler().Start()
	base.RegisterAtExit(base.GlobalScheduler.Stop)

	var mw MainWin
	mw.app = app.NewWithID(sidTheme.AppTitle) // Must Set First
	mw.app.SetIcon(sidTheme.ResourceAppIcon)

	globalWin = &mw
	globalConfig = common.NewConfig()

	// Status Bar
	mw.sb = newStatusBar()

	globalLogWriter = common.NewLogWriter(common.LogWriterOption{
		CacheCapacity: 256,
		LogPath:       mw.app.Storage().RootURI().Path(),
		OnMessage: func(s string) {
			mw.sb.setMessage(s)
		},
	})

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(globalLogWriter)

	mw.win = mw.app.NewWindow(sidTheme.AppTitle)
	mw.win.SetPadded(false)

	preTheme, _ := globalConfig.Theme.Get()
	switch preTheme {
	case "__DARK__":
		mw.app.Settings().SetTheme(sidTheme.DarkTheme{})
	case "__LIGHT__":
		mw.app.Settings().SetTheme(sidTheme.LightTheme{})
	}

	// Main Menu
	mw.mm = newMainMenu()
	mw.win.SetMainMenu(mw.mm.MainMenu)

	// Tool Bar
	mw.tb = newToolBar()

	// Fun Toys
	mw.toys = newToys()

	// Main App Tabs
	mw.at = newAppContainer()

	content := container.NewBorder(mw.tb.toolbar, mw.sb.widget, nil, mw.toys.widget, mw.at)
	mw.win.SetContent(content)
	mw.win.Resize(fyne.NewSize(800, 600))

	mw.win.SetMaster()
	mw.win.CenterOnScreen()

	mw.win.SetCloseIntercept(mw.quitHandle)
	return globalWin
}

func (mw *MainWin) Run() {
	defer func() {
		//_ = os.Unsetenv("FYNE_FONT")
		//_ = os.Unsetenv("FYNE_FONT_MONOSPACE")

		base.Exit()
	}()

	go func() {
		time.Sleep(1 * time.Second)
		mw.mm.resetMenuStatAfterMainWindowShow()
	}()

	app, err := mw.at.openDefaultApp()
	if err != nil {
		printErr(fmt.Errorf(sidTheme.RunAppFailedFormat, app, err))
	}
	log.Print(sidTheme.WelComeMsg)

	mw.win.ShowAndRun()
}

func (mw *MainWin) quitHandle() {
	// TODO. now, close window directly
	if true {
		mw.win.Close()
		return
	}
	d := dialog.NewConfirm(sidTheme.QuitAppTitle, sidTheme.QuitAppMsg, func(b bool) {
		if b {
			mw.win.Close()
		}
	}, mw.win)

	d.SetDismissText(sidTheme.DismissText)
	d.SetConfirmText(sidTheme.ConfirmText)
	d.Show()
}
