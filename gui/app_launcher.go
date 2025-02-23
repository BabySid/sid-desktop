package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	"path/filepath"
	"sid-desktop/common/apps"
	"sid-desktop/storage"
	"sid-desktop/theme"
	"time"
)

var _ appInterface = (*appLauncher)(nil)

type appLauncher struct {
	appAdapter
	searchEntry *widget.Entry
	explorer    *widget.Select
	config      *widget.Button
	configWin   fyne.Window
	appHeader   fyne.CanvasObject
	appList     *widget.List
	appBinding  binding.UntypedList
	appHistory  *apps.AppList
	appCache    *apps.AppList
}

func (al *appLauncher) LazyInit() error {
	err := storage.GetAppLauncherDB().Open(globalWin.app.Storage().RootURI().Path())
	if err != nil {
		return err
	}
	gobase.RegisterAtExit(storage.GetAppLauncherDB().Close)

	al.searchEntry = widget.NewEntry()
	al.searchEntry.SetPlaceHolder(theme.AppLauncherSearchPlaceHolder)
	al.searchEntry.OnChanged = al.searchApp
	al.searchEntry.OnSubmitted = al.execCommand

	al.explorer = widget.NewSelect(al.InitExplorerPath(), al.openExplorer)
	al.explorer.PlaceHolder = theme.AppLauncherExplorerText

	al.config = widget.NewButtonWithIcon(theme.AppLauncherConfigBtnText, theme.ResourceConfIndexIcon, al.openConfig)

	al.appBinding = binding.NewUntypedList()
	al.createAppList()

	al.tabItem = container.NewTabItemWithIcon(theme.AppLauncherName, theme.ResourceLauncherIcon, nil)
	al.tabItem.Content = container.NewBorder(
		container.NewGridWithColumns(2,
			al.searchEntry,
			container.NewHBox(layout.NewSpacer(), al.explorer, al.config)), nil, nil, nil,
		container.NewBorder(al.appHeader, nil, nil, nil, al.appList),
	)

	go al.initDB()

	return nil
}

func (al *appLauncher) GetAppName() string {
	return theme.AppLauncherName
}

func (al *appLauncher) OnClose() bool {
	if al.configWin == nil {
		return true
	}

	dialog.ShowInformation(theme.CannotCloseTitle, theme.AppLauncherCannotCloseMsg, globalWin.win)
	return false
}

func (al *appLauncher) createAppList() {
	// App List Header
	al.appHeader = container.NewGridWithColumns(3,
		widget.NewLabelWithStyle(theme.AppLauncherAppListHeader1, fyne.TextAlignLeading, fyne.TextStyle{}),
		widget.NewLabelWithStyle(theme.AppLauncherAppListHeader2, fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{}))

	// App Data
	al.appList = widget.NewListWithData(
		al.appBinding,
		func() fyne.CanvasObject {
			appName := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{})
			return container.NewGridWithColumns(3,
				container.NewHBox(widget.NewIcon(theme.ResourceDefAppIcon), appName, layout.NewSpacer()),
				widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
				container.NewHBox(
					layout.NewSpacer(),
					widget.NewButtonWithIcon("", theme.ResourceOpenDirIcon, nil),
					widget.NewButtonWithIcon("", theme.ResourceRunIcon, nil)),
			)
		},
		func(data binding.DataItem, item fyne.CanvasObject) {
			o, _ := data.(binding.Untyped).Get()
			app := o.(apps.AppInfo)
			if app.Icon != nil && len(app.Icon) > 0 {
				item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Icon).SetResource(
					fyne.NewStaticResource(app.AppName, app.Icon))
			} else {
				item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Icon).SetResource(
					theme.ResourceDefAppIcon)
			}
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(app.AppName)
			if app.AccessTime <= 0 {
				item.(*fyne.Container).Objects[1].(*widget.Label).SetText("-")
			} else {
				item.(*fyne.Container).Objects[1].(*widget.Label).SetText(gobase.FormatTimeStamp(app.AccessTime))
			}

			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[1].(*widget.Button).SetText(theme.AppLauncherAppListOp1)
			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
				err := gobase.ExecExplorer([]string{filepath.Dir(app.FullPath)})
				if err != nil {
					printErr(fmt.Errorf(theme.OpenAppLocationFailedFormat, app.FullPath, err))
				}
			}
			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[2].(*widget.Button).SetText(theme.AppLauncherAppListOp2)
			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[2].(*widget.Button).OnTapped = func() {
				app.AccessTime = time.Now().Unix()
				go func() {
					err := storage.GetAppLauncherDB().UpdateAppInfo(app)
					if err != nil {
						printErr(fmt.Errorf(theme.UpdateAppIndexFailedFormat, app.FullPath, err))
					}
					al.appCache.UpdateAppInfo(app)
					al.appHistory, err = storage.GetAppLauncherDB().LoadAppHistory()
					if err != nil {
						printErr(fmt.Errorf(theme.RunAppIndexFailedFormat, err))
					}
				}()
				err := app.Exec()
				if err != nil {
					printErr(fmt.Errorf(theme.RunAppFailedFormat, app.FullPath, err))
				}
			}
		},
	)
}

func (al *appLauncher) InitExplorerPath() []string {
	path := gobase.GetDiskPartitions()

	path = append(path, theme.AppLauncherExplorerSidPathName)
	return path
}

func (al *appLauncher) loadAppInfoFromDB() {
	var err error
	al.appCache, err = storage.GetAppLauncherDB().LoadAppIndex()
	if err != nil {
		printErr(fmt.Errorf(theme.RunAppIndexFailedFormat, err))
	}

	al.appHistory, err = storage.GetAppLauncherDB().LoadAppHistory()
	if err != nil {
		printErr(fmt.Errorf(theme.RunAppIndexFailedFormat, err))
	}
}

func (al *appLauncher) openConfig() {
	if al.configWin == nil {
		al.configWin = newAppLauncherConfig(al).win
		al.configWin.Show()
		al.configWin.SetOnClosed(func() {
			al.configWin = nil
		})
	} else {
		al.configWin.RequestFocus()
	}
}

func (al *appLauncher) searchApp(name string) {
	if name == "" {
		// show history
		if al.appHistory != nil {
			_ = al.appBinding.Set(al.appHistory.AsInterfaceArray())
		}
	} else {
		if al.appCache != nil {
			rs := al.appCache.Find(name)
			_ = al.appBinding.Set(rs.AsInterfaceArray())
		}
	}
}

func (al *appLauncher) execCommand(cmd string) {
	err := gobase.ExecApp(cmd)
	if err != nil {
		printErr(fmt.Errorf(theme.RunCommandFailedFormat, cmd, err))
	}
}

func (al *appLauncher) openExplorer(dir string) {
	path := make([]string, 1, 1)
	if dir == theme.AppLauncherExplorerSidPathName {
		path[0], _ = filepath.Abs(globalWin.app.Storage().RootURI().Path())
	} else {
		path[0] = dir
	}
	err := gobase.ExecExplorer(path)
	if err != nil {
		printErr(fmt.Errorf(theme.RunExplorerFailedFormat, dir, err))
	}
}

func (al *appLauncher) initDB() {
	need, err := storage.GetAppLauncherDB().NeedInit()
	if err != nil {
		printErr(fmt.Errorf(theme.RunAppIndexFailedFormat, err))
		return
	}
	if need {
		dialog.ShowConfirm(theme.AppLauncherNeedInitTitle, theme.AppLauncherNeedInitMsg, func(b bool) {
			if b {
				al.openConfig()
			}
		}, globalWin.win)
	} else {
		al.loadAppInfoFromDB()
		if al.appHistory != nil {
			_ = al.appBinding.Set(al.appHistory.AsInterfaceArray())
		}
	}
}

func (al *appLauncher) ShortCut() fyne.Shortcut {
	return &desktop.CustomShortcut{KeyName: fyne.Key2, Modifier: fyne.KeyModifierAlt}
}

func (al *appLauncher) Icon() fyne.Resource {
	return theme.ResourceLauncherIcon
}
