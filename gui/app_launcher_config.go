package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	"sid-desktop/common"
	"sid-desktop/common/apps"
	"sid-desktop/storage"
	theme2 "sid-desktop/theme"
	"time"
)

type appLauncherConfig struct {
	al *appLauncher

	ok      *widget.Button
	dismiss *widget.Button

	indexBuildStatus int
	indexBuildLabel  *widget.Label

	pathBinding binding.StringList

	win fyne.Window
}

const (
	indexBuildReady = iota
	indexBuilding
	indexBuildFinished
)

func newAppLauncherConfig(launcher *appLauncher) *appLauncherConfig {
	var alc appLauncherConfig

	alc.al = launcher

	alc.indexBuildStatus = indexBuildReady
	alc.indexBuildLabel = widget.NewLabel("")
	alc.indexBuildLabel.Hide()

	alc.win = fyne.CurrentApp().NewWindow(theme2.AppLauncherConfigTitle)

	alc.pathBinding = binding.NewStringList()
	addFolder := widget.NewButtonWithIcon(theme2.AppLauncherConfigAddDirBtn, theme2.ResourceAddDirIcon, func() {
		fo := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil {
				_ = alc.pathBinding.Append(uri.Path())
			}
		}, alc.win)
		fo.Show()
	})

	fileFilter := widget.NewLabel("*.lnk;*.exe")
	top := container.NewHBox(widget.NewLabel(theme2.AppLauncherConfigFileFilter), fileFilter, layout.NewSpacer(), addFolder)

	common.CopyBindingStringList(alc.pathBinding, globalConfig.AppLaunchAppSearchPath)

	pathCont := widget.NewListWithData(
		alc.pathBinding,
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				layout.NewSpacer(),
				widget.NewButtonWithIcon(theme2.AppLauncherConfigRmDirBtn, theme2.ResourceRmDirIcon, nil))
		},
		func(data binding.DataItem, item fyne.CanvasObject) {
			path, _ := data.(binding.String).Get()
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(path)
			item.(*fyne.Container).Objects[2].(*widget.Button).OnTapped = func() {
				pathTemp, _ := alc.pathBinding.Get()
				idx := gobase.ContainsString(pathTemp, path)
				pathTemp = append(pathTemp[:idx], pathTemp[idx+1:]...)
				_ = alc.pathBinding.Set(pathTemp)
			}
		},
	)

	alc.ok = widget.NewButtonWithIcon(theme2.AppLauncherConfigBuildBtn, theme.ConfirmIcon(), func() {
		alc.indexBuildStatus = indexBuilding
		alc.indexBuildLabel.SetText(theme2.AppLauncherConfigIndexBuilding)
		alc.indexBuildLabel.Show()
		alc.ok.Disable()
		alc.dismiss.Disable()

		go func() {
			defer func() {
				alc.indexBuildStatus = indexBuildFinished
				alc.ok.Enable()
				alc.dismiss.Enable()
			}()
			common.CopyBindingStringList(globalConfig.AppLaunchAppSearchPath, alc.pathBinding)

			time.Sleep(200 * time.Millisecond)
			err := storage.GetAppLauncherDB().Init()
			if err != nil {
				printErr(fmt.Errorf(theme2.RunAppIndexFailedFormat, err))
				return
			}
			alc.indexBuildLabel.SetText(theme2.AppLauncherConfigStartScanApp)
			time.Sleep(200 * time.Millisecond)

			path, _ := globalConfig.AppLaunchAppSearchPath.Get()
			appFound, err := apps.InitApps(path)
			if err != nil {
				printErr(fmt.Errorf(theme2.RunAppIndexFailedFormat, err))
				return
			}

			err = storage.GetAppLauncherDB().AddAppToIndex(appFound)
			if err != nil {
				printErr(fmt.Errorf(theme2.RunAppIndexFailedFormat, err))
				return
			}
			alc.notifyAppLauncherBuildIndexFinished()

			alc.indexBuildLabel.SetText(fmt.Sprintf(theme2.AppLauncherConfigFinishScanAppFormat, appFound.Len()))
		}()
	})
	alc.dismiss = widget.NewButtonWithIcon(theme2.DismissText, theme.CancelIcon(), alc.closeHandle)

	alc.win.SetContent(container.NewBorder(top,
		container.NewVBox(
			alc.indexBuildLabel,
			container.NewHBox(layout.NewSpacer(), alc.dismiss, alc.ok),
		),
		nil, nil,
		pathCont))

	alc.win.SetCloseIntercept(alc.closeHandle)
	alc.win.Resize(fyne.NewSize(600, 300))
	alc.win.CenterOnScreen()
	return &alc
}

func (alc *appLauncherConfig) closeHandle() {
	if alc.indexBuildStatus == indexBuilding {
		dialog.ShowInformation(theme2.CannotCloseTitle, theme2.AppLauncherConfigCannotCloseMsg, globalWin.win)
	} else {
		alc.win.Close()
	}
}

func (alc *appLauncherConfig) notifyAppLauncherBuildIndexFinished() {
	alc.al.loadAppInfoFromDB()
}
