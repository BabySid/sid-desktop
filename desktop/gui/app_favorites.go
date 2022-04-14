package gui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	"sid-desktop/desktop/common"
	"sid-desktop/desktop/storage"
	sidTheme "sid-desktop/desktop/theme"
	"strings"
	"time"
)

var _ appInterface = (*appFavorites)(nil)

type appFavorites struct {
	appAdapter
	searchEntry  *widget.Entry
	importFavor  *widget.Button
	newFavor     *widget.Button
	exportFavor  *widget.Button
	favorHeader  *widget.List
	favorList    *widget.List
	favorBinding binding.UntypedList
	favorCache   *common.FavoritesList
}

func (af *appFavorites) LazyInit() error {
	err := storage.GetAppFavoritesDB().Open(globalWin.app.Storage().RootURI().Path())
	if err != nil {
		return err
	}
	gobase.RegisterAtExit(storage.GetAppFavoritesDB().Close)

	af.searchEntry = widget.NewEntry()
	af.searchEntry.SetPlaceHolder(sidTheme.AppFavoritesSearchText)
	af.searchEntry.OnChanged = af.searchFavor

	af.newFavor = widget.NewButtonWithIcon(sidTheme.AppFavoritesAddFavorBtnText, sidTheme.ResourceAddIcon, af.addFavor)
	af.importFavor = widget.NewButtonWithIcon(sidTheme.AppFavoritesImportBtnText, sidTheme.ResourceImportIcon, af.importFavors)
	af.exportFavor = widget.NewButtonWithIcon(sidTheme.AppFavoritesExportBtnText, sidTheme.ResourceExportIcon, af.exportFavors)

	af.favorBinding = binding.NewUntypedList()
	af.createFavorList()

	af.tabItem = container.NewTabItemWithIcon(sidTheme.AppFavoritesName, sidTheme.ResourceFavoritesIcon, nil)
	af.tabItem.Content = container.NewBorder(
		container.NewGridWithColumns(2,
			af.searchEntry,
			container.NewHBox(layout.NewSpacer(), af.importFavor, af.exportFavor, af.newFavor)), nil, nil, nil,
		container.NewBorder(af.favorHeader, nil, nil, nil, af.favorList),
	)

	go af.initDB()

	return nil
}

func (af *appFavorites) GetAppName() string {
	return sidTheme.AppFavoritesName
}

func (af *appFavorites) exportFavors() {
	d := dialog.NewFileSave(func(closer fyne.URIWriteCloser, err error) {
		if err != nil {
			printErr(fmt.Errorf(sidTheme.ExportFavoritesFailedFormat, err))
			return
		}
		if closer != nil {
			defer closer.Close()
			favors := af.favorCache.GetFavorites()

			data := ""
			for _, favor := range favors {
				t, _ := json.Marshal(favor)
				if data != "" {
					data += "\n"
				}
				data += string(t)
			}

			_, err := closer.Write([]byte(data))
			if err != nil {
				printErr(fmt.Errorf(sidTheme.ExportFavoritesFailedFormat, err))
			}
		}
	}, globalWin.win)
	d.Show()
}

func (af *appFavorites) importFavors() {
	d := dialog.NewFileOpen(func(closer fyne.URIReadCloser, err error) {
		if err != nil {
			printErr(fmt.Errorf(sidTheme.ImportFavoritesFailedFormat, err))
			return
		}

		if closer != nil {
			defer closer.Close()

			data, err := common.ReadURI(closer)
			if err != nil {
				printErr(fmt.Errorf(sidTheme.ImportFavoritesFailedFormat, err))
				return
			}

			favorBytes := bytes.Split(data, []byte("\n"))

			favors := common.NewFavoritesList()
			for _, item := range favorBytes {
				var fav common.Favorites

				err = json.Unmarshal(item, &fav)
				if err != nil {
					printErr(fmt.Errorf(sidTheme.ImportFavoritesFailedFormat, err))
					return
				}

				favors.Append(fav)
			}

			err = storage.GetAppFavoritesDB().AddFavoritesList(favors)
			if err != nil {
				printErr(fmt.Errorf(sidTheme.ImportFavoritesFailedFormat, err))
			}

			af.reloadFavorList()
		}
	}, globalWin.win)
	d.Show()
}

func (af *appFavorites) createFavorList() {
	// Favor List Header
	af.favorHeader = widget.NewList(
		func() int {
			return 1
		},
		func() fyne.CanvasObject {
			return container.NewGridWithColumns(3,
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
				widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{}))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(sidTheme.AppFavoritesFavorListHeader1)
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(sidTheme.AppFavoritesFavorListHeader2)
			item.(*fyne.Container).Objects[2].(*widget.Label).SetText("")
		},
	)

	// Favor Data
	af.favorList = widget.NewListWithData(
		af.favorBinding,
		func() fyne.CanvasObject {
			return container.NewGridWithColumns(3,
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
				container.NewHBox(
					layout.NewSpacer(),
					widget.NewButtonWithIcon(sidTheme.AppFavoritesFavorListOp1, sidTheme.ResourceEditIcon, nil),
					widget.NewButtonWithIcon(sidTheme.AppFavoritesFavorListOp2, sidTheme.ResourceOpenUrlIcon, nil)),
			)
		},
		func(data binding.DataItem, item fyne.CanvasObject) {
			o, _ := data.(binding.Untyped).Get()
			favor := o.(common.Favorites)

			name := gobase.CutUTF8(favor.Name, 0, 16, "...")
			name += "(" + strings.Join(favor.Tags, common.FavorTagSep) + ")"
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(name)

			localUrl := gobase.CutUTF8(favor.Url, 0, 64, "...")
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(localUrl)

			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[1].(*widget.Button).SetText(sidTheme.AppFavoritesFavorListOp1)
			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
				af.editOneFavor(favor)
			}
			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[2].(*widget.Button).SetText(sidTheme.AppFavoritesFavorListOp2)
			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[2].(*widget.Button).OnTapped = func() {
				err := globalWin.app.OpenURL(common.ParseURL(favor.Url))
				if err != nil {
					printErr(fmt.Errorf(sidTheme.OpenFavoritesFailedFormat, err))
				}
			}
		},
	)
}

func (af *appFavorites) searchFavor(name string) {
	if name == "" {
		// show all favorites
		if af.favorCache != nil {
			_ = af.favorBinding.Set(af.favorCache.AsInterfaceArray())
		}
	} else {
		if af.favorCache != nil {
			rs := af.favorCache.Find(name)
			_ = af.favorBinding.Set(rs.AsInterfaceArray())
		}
	}
}

func (af *appFavorites) editOneFavor(favor common.Favorites) {
	af.showFavorDialog(&favor)
}

func (af *appFavorites) addFavor() {
	af.showFavorDialog(nil)
}

func (af *appFavorites) showFavorDialog(favor *common.Favorites) {
	url := widget.NewEntry()
	url.Validator = validation.NewRegexp(
		`(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`,
		"please input right URL")

	url.SetPlaceHolder(sidTheme.AppFavoritesAddFavorUrlPlaceHolder)

	name := widget.NewEntry()
	name.Validator = validation.NewRegexp(`\S+`, sidTheme.AppFavoritesAddFavorName+" must not be empty")

	tags := widget.NewEntry()
	expand := widget.NewButtonWithIcon(sidTheme.AppFavoritesAddFavorExpand, sidTheme.ResourceExpandDownIcon, nil)

	url.OnSubmitted = func(s string) {
		tl, err := gobase.GetWebPageTitle(s)
		if err != nil {
			printErr(fmt.Errorf(sidTheme.WebPageProcessErrorFormat, err))
		}

		if tl != "" {
			name.SetText(tl)
		}
	}

	var rmBtn *widget.Button
	if favor != nil {
		rmBtn = widget.NewButtonWithIcon(sidTheme.AppFavoritesRmFavorBtnText, sidTheme.ResourceRmIcon, nil)
	}

	tagArray := binding.NewStringList()
	tagList := widget.NewListWithData(
		tagArray,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(item binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(item.(binding.String))
		},
	)
	tags.OnChanged = func(s string) {
		arr := strings.Split(s, common.FavorTagSep)
		_ = tagArray.Set(arr)
	}

	title := sidTheme.AppFavoritesAddFavorTitle
	if favor != nil { // edit or remove
		url.SetText(favor.Url)
		name.SetText(favor.Name)
		tags.SetText(strings.Join(favor.Tags, common.FavorTagSep))
		title = sidTheme.AppFavoritesEditFavorTitle
	}

	tagList.Hide()

	var opContainer *fyne.Container
	if favor != nil {
		opContainer = container.NewHBox(expand, layout.NewSpacer(), rmBtn)
	} else {
		opContainer = container.NewHBox(expand, layout.NewSpacer())
	}
	cont := container.NewBorder(
		container.NewVBox(
			widget.NewForm(
				widget.NewFormItem(sidTheme.AppFavoritesAddFavorUrl, url),
				widget.NewFormItem(sidTheme.AppFavoritesAddFavorName, name),
				widget.NewFormItem(sidTheme.AppFavoritesAddFavorTags, tags),
			),
			opContainer,
		),
		nil, nil, nil,
		tagList,
	)

	win := dialog.NewCustomConfirm(
		title, sidTheme.ConfirmText, sidTheme.DismissText,
		cont, func(b bool) {
			if b {
				t, _ := tagArray.Get()

				var tempFavor common.Favorites
				if favor != nil {
					tempFavor = *favor
				}
				tempFavor.Name = name.Text
				tempFavor.Url = url.Text
				tempFavor.Tags = t
				tempFavor.CreateTime = time.Now().Unix()

				if favor != nil {
					err := storage.GetAppFavoritesDB().UpdateFavorites(tempFavor)
					if err != nil {
						printErr(fmt.Errorf(sidTheme.ProcessFavoritesFailedFormat, err))
					}
				} else {
					err := storage.GetAppFavoritesDB().AddFavorites(tempFavor)
					if err != nil {
						printErr(fmt.Errorf(sidTheme.ProcessFavoritesFailedFormat, err))
					}
				}

				af.reloadFavorList()
			}
		},
		globalWin.win)

	expand.OnTapped = func() {
		if tagList.Visible() {
			tagList.Hide()
			win.Resize(fyne.NewSize(500, 300))

			expand.SetText(sidTheme.AppFavoritesAddFavorExpand)
			expand.SetIcon(sidTheme.ResourceExpandDownIcon)
		} else {
			tagList.Show()
			win.Resize(fyne.NewSize(500, 500))

			expand.SetText(sidTheme.AppFavoritesAddFavorShrink)
			expand.SetIcon(sidTheme.ResourceExpandUpIcon)
		}
	}

	if rmBtn != nil {
		rmBtn.OnTapped = func() {
			_ = storage.GetAppFavoritesDB().RmFavorites(*favor)
			af.reloadFavorList()
			win.Hide()
		}
	}

	win.Resize(fyne.NewSize(500, 300))

	win.Show()
}

func (af *appFavorites) reloadFavorList() {
	af.loadFavoritesFromDB()
	if af.favorCache != nil {
		_ = af.favorBinding.Set(af.favorCache.AsInterfaceArray())
	}
}

func (af *appFavorites) loadFavoritesFromDB() {
	var err error
	af.favorCache, err = storage.GetAppFavoritesDB().LoadFavorites()
	if err != nil {
		printErr(fmt.Errorf(sidTheme.ProcessFavoritesFailedFormat, err))
	}
}

func (af *appFavorites) initDB() {
	need, err := storage.GetAppFavoritesDB().NeedInit()
	if err != nil {
		printErr(fmt.Errorf(sidTheme.ProcessFavoritesFailedFormat, err))
		return
	}

	if need {
		err = storage.GetAppFavoritesDB().Init()
		if err != nil {
			printErr(fmt.Errorf(sidTheme.ProcessFavoritesFailedFormat, err))
			return
		}
	} else {
		af.reloadFavorList()
	}
}
