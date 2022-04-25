package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"sid-desktop/desktop/common"
	sidTheme "sid-desktop/desktop/theme"
)

var _ appInterface = (*appWelcome)(nil)

type appWelcome struct {
	appAdapter
}

func (a *appWelcome) LazyInit() error {
	logo := canvas.NewImageFromResource(sidTheme.ResourceSidLogo)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(362*0.8, 192*0.8))

	a.tabItem = container.NewTabItemWithIcon(sidTheme.AppWelcomeName, sidTheme.ResourceWelIcon, nil)

	wel := widget.NewRichTextFromMarkdown("# " + sidTheme.WelComeMsg)
	for i := range wel.Segments {
		if seg, ok := wel.Segments[i].(*widget.TextSegment); ok {
			seg.Style.Alignment = fyne.TextAlignCenter
		}
	}

	shortCuts := widget.NewForm()
	for _, myApp := range appRegister {
		shortCuts.Append(myApp.GetAppName(), widget.NewLabelWithStyle(common.ShortCutName(myApp.ShortCut()), fyne.TextAlignCenter, fyne.TextStyle{}))
	}
	a.tabItem.Content = container.NewCenter(
		container.NewVBox(
			wel,
			logo,
			shortCuts,
		))

	return nil
}

func (a *appWelcome) GetAppName() string {
	return sidTheme.AppWelcomeName
}

func (a *appWelcome) OpenDefault() bool {
	return true
}

func (a *appWelcome) ShortCut() fyne.Shortcut {
	return &desktop.CustomShortcut{KeyName: fyne.Key1, Modifier: desktop.AltModifier}
}
