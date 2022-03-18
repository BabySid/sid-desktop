package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"net/url"
	sidTheme "sid-desktop/desktop/theme"
)

var _ appInterface = (*appWelcome)(nil)

type appWelcome struct {
	tabItem *container.TabItem
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

	a.tabItem.Content = container.NewCenter(
		container.NewVBox(
			wel,
			logo,
			container.NewHBox(
				widget.NewHyperlink("百度", parseURL("https://fyne.io/")),
				widget.NewLabel("-"),
				widget.NewHyperlink("documentation", parseURL("https://developer.fyne.io/")),
				widget.NewLabel("-"),
				widget.NewHyperlink("sponsor", parseURL("https://fyne.io/sponsor/")),
			),
		))
	return nil
}

func (a *appWelcome) GetTabItem() *container.TabItem {
	return a.tabItem
}

func (a *appWelcome) GetAppName() string {
	return sidTheme.AppWelcomeName
}

func (a *appWelcome) OpenDefault() bool {
	return true
}

func (a *appWelcome) OnClose() bool {
	return true
}

func parseURL(urlStr string) *url.URL {
	link, _ := url.Parse(urlStr)

	return link
}
