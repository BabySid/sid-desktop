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
	sidTheme "sid-desktop/theme"
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

	alc.win = fyne.CurrentApp().NewWindow(sidTheme.AppLauncherConfigTitle)

	alc.pathBinding = binding.NewStringList()
	addFolder := widget.NewButtonWithIcon(sidTheme.AppLauncherConfigAddDirBtn, sidTheme.ResourceAddDirIcon, func() {
		fo := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if uri != nil {
				_ = alc.pathBinding.Append(uri.Path())
			}
		}, alc.win)
		fo.Show()
	})

	fileFilter := widget.NewLabel("*.lnk;*.exe")
	top := container.NewHBox(widget.NewLabel(sidTheme.AppLauncherConfigFileFilter), fileFilter, layout.NewSpacer(), addFolder)

	common.CopyBindingStringList(alc.pathBinding, globalConfig.AppLaunchAppSearchPath)

	pathCont := widget.NewListWithData(
		alc.pathBinding,
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				layout.NewSpacer(),
				widget.NewButtonWithIcon(sidTheme.AppLauncherConfigRmDirBtn, sidTheme.ResourceRmDirIcon, nil))
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

	alc.ok = widget.NewButtonWithIcon(sidTheme.AppLauncherConfigBuildBtn, theme.ConfirmIcon(), func() {
		alc.indexBuildStatus = indexBuilding
		alc.indexBuildLabel.SetText(sidTheme.AppLauncherConfigIndexBuilding)
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
				printErr(fmt.Errorf(sidTheme.RunAppIndexFailedFormat, err))
				return
			}
			alc.indexBuildLabel.SetText(sidTheme.AppLauncherConfigStartScanApp)
			time.Sleep(200 * time.Millisecond)

			path, _ := globalConfig.AppLaunchAppSearchPath.Get()
			appFound, err := apps.InitApps(path)
			if err != nil {
				printErr(fmt.Errorf(sidTheme.RunAppIndexFailedFormat, err))
				return
			}

			err = storage.GetAppLauncherDB().AddAppToIndex(appFound)
			if err != nil {
				printErr(fmt.Errorf(sidTheme.RunAppIndexFailedFormat, err))
				return
			}
			alc.notifyAppLauncherBuildIndexFinished()

			alc.indexBuildLabel.SetText(fmt.Sprintf(sidTheme.AppLauncherConfigFinishScanAppFormat, appFound.Len()))
		}()
	})
	alc.dismiss = widget.NewButtonWithIcon(sidTheme.DismissText, theme.CancelIcon(), alc.closeHandle)

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
		dialog.ShowInformation(sidTheme.CannotCloseTitle, sidTheme.AppLauncherConfigCannotCloseMsg, globalWin.win)
	} else {
		alc.win.Close()
	}
}

func (alc *appLauncherConfig) notifyAppLauncherBuildIndexFinished() {
	alc.al.loadAppInfoFromDB()
}
