package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/BabySid/gobase"
	"github.com/vicanso/go-charts/v2"
	"golang.design/x/hotkey"
	"log"
	"os"
	"sid-desktop/common"
	"sid-desktop/theme"
)

type winStatus struct {
	shown bool
}
type MainWin struct {
	app   fyne.App
	win   fyne.Window
	wStat winStatus
	mm    *mainMenu
	tb    *toolBar
	toys  *toys
	sb    *statusBar
	at    *appContainer
}

func init() {
	// set env to support chinese
	//_ = os.Setenv("FYNE_FONT", "./resource/Microsoft-YaHei.ttf")
	//_ = os.Setenv("FYNE_FONT_MONOSPACE", "./resource/Microsoft-YaHei.ttf")
	// _ = os.Setenv("FYNE_SCALE", "0.8")
	_ = charts.InstallFont("yahei", theme.FontMicrosoftYaHeiTtf.StaticContent)
	font, _ := charts.GetFont("yahei")
	charts.SetDefaultFont(font)
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
	globalLogWriter *common.LogWriter

	tray *sysTray
)

func NewMainWin() *MainWin {
	gobase.NewScheduler().Start()
	gobase.RegisterAtExit(gobase.GlobalScheduler.Stop)

	var mw MainWin
	mw.app = app.NewWithID(theme.AppTitle) // Must Set First
	mw.app.SetIcon(theme.ResourceAppIcon)

	mw.win = mw.app.NewWindow(theme.AppTitle)
	mw.win.SetPadded(false)

	globalWin = &mw

	globalLogWriter = common.NewLogWriter(common.LogWriterOption{
		CacheCapacity: 256,
		LogPath:       mw.app.Storage().RootURI().Path(),
		OnMessage: func(s string) {
			if mw.sb == nil {
				return
			}
			mw.sb.setMessage(s)
		},
	})

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(globalLogWriter)

	preTheme, _ := common.GetConfig().Theme.Get()
	switch preTheme {
	case "__DARK__":
		mw.app.Settings().SetTheme(theme.DarkTheme)
	case "__LIGHT__":
		mw.app.Settings().SetTheme(theme.LightTheme)
	}

	// todo setup splash window

	// Main Menu
	mw.mm = newMainMenu()
	mw.win.SetMainMenu(mw.mm.MainMenu)

	// Tool Bar
	mw.tb = newToolBar()

	// Fun Toys
	mw.toys = newToys()

	// Main App Tabs
	mw.at = newAppContainer()

	// Status Bar
	mw.sb = newStatusBar()

	content := container.NewBorder(mw.tb.toolbar, mw.sb.widget, nil, mw.toys.widget, mw.at)
	mw.win.SetContent(content)
	mw.win.Resize(fyne.NewSize(1024, 768))

	mw.win.SetMaster()
	mw.win.CenterOnScreen()

	mw.win.SetCloseIntercept(mw.quitHandle)
	mw.wStat.shown = true

	if _, ok := mw.app.(desktop.App); ok {
		tray = newSysTray()
	}

	mw.registerShortCut()

	return globalWin
}

func (mw *MainWin) Run() {
	defer func() {
		//_ = os.Unsetenv("FYNE_FONT")
		//_ = os.Unsetenv("FYNE_FONT_MONOSPACE")
		_ = os.Unsetenv("FYNE_SCALE")

		gobase.Exit()
	}()

	appName, err := mw.at.openDefaultApp()
	if err != nil {
		printErr(fmt.Errorf(theme.RunAppFailedFormat, appName, err))
	}
	log.Print(theme.WelComeMsg)

	// set up systray
	if tray != nil {
		tray.run()
	}

	mw.win.ShowAndRun()
}

func (mw *MainWin) quitHandle() {
	hide, _ := common.GetConfig().HideWhenQuit.Get()
	if hide {
		mw.hideWin()
		return
	}
	mw.closeWin()
	//d := dialog.NewConfirm(sidTheme.QuitAppTitle, sidTheme.QuitAppMsg, func(b bool) {
	//	if b {
	//		mw.win.Close()
	//	}
	//}, mw.win)
	//
	//d.SetDismissText(sidTheme.DismissText)
	//d.SetConfirmText(sidTheme.ConfirmText)
	//d.Show()
}

func (mw *MainWin) closeWin() {
	if tray != nil {
		tray.Quit()
	}
	// TODO. now, close window directly
	mw.win.Close()
}

func (mw *MainWin) showWin() {
	mw.win.Show()
	mw.wStat.shown = true

	tray.setHideMenu()
}

func (mw *MainWin) hideWin() {
	mw.win.Hide()
	mw.wStat.shown = false

	tray.setShowMenu()
}

func (mw *MainWin) registerShortCut() {
	for _, myApp := range appRegister {
		ap := myApp
		sc := ap.ShortCut()
		mw.win.Canvas().AddShortcut(sc, func(_ fyne.Shortcut) {
			err := globalWin.at.openApp(ap)
			if err != nil {
				printErr(fmt.Errorf(theme.RunAppFailedFormat, ap.GetAppName(), err))
			}
		})
	}

	go func() {
		hk := hotkey.New([]hotkey.Modifier{hotkey.ModCtrl, hotkey.ModShift}, hotkey.KeyZ)
		if err := hk.Register(); err != nil {
			printErr(fmt.Errorf(theme.InternalErrorFormat, err))
			return
		}
		defer hk.Unregister()
		for range hk.Keydown() {
			if globalWin.wStat.shown {
				globalWin.hideWin()
			} else {
				globalWin.showWin()
			}
		}
	}()
}
