package gui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	"image/color"
	"sid-desktop/common"
	"sid-desktop/storage"
	"sid-desktop/theme"
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
	favorHeader  fyne.CanvasObject
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
	af.searchEntry.SetPlaceHolder(theme.AppFavoritesSearchText)
	af.searchEntry.OnChanged = af.searchFavor

	af.newFavor = widget.NewButtonWithIcon(theme.AppFavoritesAddFavorBtnText, theme.ResourceAddIcon, af.addFavor)
	af.importFavor = widget.NewButtonWithIcon(theme.AppFavoritesImportBtnText, theme.ResourceImportIcon, af.importFavors)
	af.exportFavor = widget.NewButtonWithIcon(theme.AppFavoritesExportBtnText, theme.ResourceExportIcon, af.exportFavors)

	af.favorBinding = binding.NewUntypedList()
	af.createFavorList()

	af.tabItem = container.NewTabItemWithIcon(theme.AppFavoritesName, theme.ResourceFavoritesIcon, nil)
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
	return theme.AppFavoritesName
}

func (af *appFavorites) exportFavors() {
	d := dialog.NewFileSave(func(closer fyne.URIWriteCloser, err error) {
		if err != nil {
			printErr(fmt.Errorf(theme.ExportFavoritesFailedFormat, err))
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
				printErr(fmt.Errorf(theme.ExportFavoritesFailedFormat, err))
			}
		}
	}, globalWin.win)
	d.Show()
}

func (af *appFavorites) importFavors() {
	d := dialog.NewFileOpen(func(closer fyne.URIReadCloser, err error) {
		if err != nil {
			printErr(fmt.Errorf(theme.ImportFavoritesFailedFormat, err))
			return
		}

		if closer != nil {
			defer closer.Close()

			data, err := common.ReadURI(closer)
			if err != nil {
				printErr(fmt.Errorf(theme.ImportFavoritesFailedFormat, err))
				return
			}

			favorBytes := bytes.Split(data, []byte("\n"))

			favors := common.NewFavoritesList()
			for _, item := range favorBytes {
				var fav common.Favorites

				err = json.Unmarshal(item, &fav)
				if err != nil {
					printErr(fmt.Errorf(theme.ImportFavoritesFailedFormat, err))
					return
				}

				favors.Append(fav)
			}

			err = storage.GetAppFavoritesDB().AddFavoritesList(favors)
			if err != nil {
				printErr(fmt.Errorf(theme.ImportFavoritesFailedFormat, err))
			}

			af.reloadFavorList()
		}
	}, globalWin.win)
	d.Show()
}

func (af *appFavorites) createFavorList() {
	// Favor List Header
	af.favorHeader = container.NewGridWithColumns(3,
		widget.NewLabelWithStyle(theme.AppFavoritesFavorListHeader1, fyne.TextAlignLeading, fyne.TextStyle{}),
		widget.NewLabelWithStyle(theme.AppFavoritesFavorListHeader2, fyne.TextAlignCenter, fyne.TextStyle{}),
		widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{}))

	// Favor Data
	af.favorList = widget.NewListWithData(
		af.favorBinding,
		func() fyne.CanvasObject {
			return container.NewGridWithColumns(3,
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
				container.NewHBox(
					layout.NewSpacer(),
					widget.NewButtonWithIcon(theme.AppFavoritesFavorListOp1, theme.ResourceEditIcon, nil),
					widget.NewButtonWithIcon(theme.AppFavoritesFavorListOp2, theme.ResourceOpenUrlIcon, nil)),
			)
		},
		func(data binding.DataItem, item fyne.CanvasObject) {
			o, _ := data.(binding.Untyped).Get()
			favor := o.(common.Favorites)

			name := gobase.CutUTF8(favor.Name, 0, 16, "...")
			name += "(" + strings.Join(favor.Tags, common.ArraySeparator) + ")"
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(name)

			localUrl := gobase.CutUTF8(favor.Url, 0, 64, "...")
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(localUrl)

			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[1].(*widget.Button).SetText(theme.AppFavoritesFavorListOp1)
			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
				af.editOneFavor(favor)
			}
			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[2].(*widget.Button).SetText(theme.AppFavoritesFavorListOp2)
			item.(*fyne.Container).Objects[2].(*fyne.Container).Objects[2].(*widget.Button).OnTapped = func() {
				err := globalWin.app.OpenURL(common.ParseURL(favor.Url))
				if err != nil {
					printErr(fmt.Errorf(theme.OpenFavoritesFailedFormat, err))
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

	url.SetPlaceHolder(theme.AppFavoritesAddFavorUrlPlaceHolder)

	name := widget.NewEntry()
	name.Validator = validation.NewRegexp(`\S+`, theme.AppFavoritesAddFavorName+" must not be empty")

	tags := widget.NewEntry()
	expand := widget.NewButtonWithIcon(theme.AppFavoritesAddFavorExpand, theme.ResourceExpandDownIcon, nil)

	url.OnSubmitted = func(s string) {
		tl, err := gobase.GetWebPageTitle(s)
		if err != nil {
			printErr(fmt.Errorf(theme.WebPageProcessErrorFormat, err))
		}

		if tl != "" {
			name.SetText(tl)
		}
	}

	var rmBtn *widget.Button
	if favor != nil {
		rmBtn = widget.NewButtonWithIcon(theme.AppFavoritesRmFavorBtnText, theme.ResourceRmIcon, nil)
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
		arr := strings.Split(s, common.ArraySeparator)
		_ = tagArray.Set(arr)
	}

	title := theme.AppFavoritesAddFavorTitle
	if favor != nil { // edit or remove
		url.SetText(favor.Url)
		name.SetText(favor.Name)
		tags.SetText(strings.Join(favor.Tags, common.ArraySeparator))
		title = theme.AppFavoritesEditFavorTitle
	}

	tagBack := canvas.NewRectangle(color.Transparent)
	tagBack.SetMinSize(fyne.NewSize(300, 50))
	tagContent := container.NewMax(tagList, tagBack)
	tagContent.Hide()

	var opContainer *fyne.Container
	if favor != nil {
		opContainer = container.NewHBox(expand, layout.NewSpacer(), rmBtn)
	} else {
		opContainer = container.NewHBox(expand, layout.NewSpacer())
	}

	win := dialog.NewForm(title, theme.ConfirmText, theme.DismissText, []*widget.FormItem{
		widget.NewFormItem(theme.AppFavoritesAddFavorUrl, url),
		widget.NewFormItem(theme.AppFavoritesAddFavorName, name),
		widget.NewFormItem(theme.AppFavoritesAddFavorTags, tags),
		widget.NewFormItem("", container.NewVBox(opContainer, tagContent)),
	}, func(b bool) {
		if !b {
			return
		}

		t, _ := tagArray.Get()

		var tempFavor common.Favorites
		if favor != nil {
			tempFavor = *favor
		}
		tempFavor.Name = name.Text
		tempFavor.Url = url.Text
		tempFavor.Tags = t
		tempFavor.CreateTime = time.Now().Unix()
		tempFavor.AccessTime = time.Now().Unix()

		if favor != nil {
			err := storage.GetAppFavoritesDB().UpdateFavorites(tempFavor)
			if err != nil {
				printErr(fmt.Errorf(theme.AppFavoritesFailedFormat, err))
			}
		} else {
			err := storage.GetAppFavoritesDB().AddFavorites(tempFavor)
			if err != nil {
				printErr(fmt.Errorf(theme.AppFavoritesFailedFormat, err))
			}
		}

		af.reloadFavorList()

	}, globalWin.win)

	expand.OnTapped = func() {
		if tagContent.Visible() {
			tagContent.Hide()
			win.Resize(fyne.NewSize(500, 300))

			expand.SetText(theme.AppFavoritesAddFavorExpand)
			expand.SetIcon(theme.ResourceExpandDownIcon)
		} else {
			tagContent.Show()
			size := tagArray.Length()
			if size >= 5 {
				size = 5
			}
			tagBack.SetMinSize(fyne.NewSize(300, common.GetItemsHeightInList(tagList, size)))
			win.Resize(fyne.NewSize(500, 500))

			expand.SetText(theme.AppFavoritesAddFavorShrink)
			expand.SetIcon(theme.ResourceExpandUpIcon)
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
		printErr(fmt.Errorf(theme.AppFavoritesFailedFormat, err))
	}
}

func (af *appFavorites) initDB() {
	need, err := storage.GetAppFavoritesDB().NeedInit()
	if err != nil {
		printErr(fmt.Errorf(theme.AppFavoritesFailedFormat, err))
		return
	}

	if need {
		err = storage.GetAppFavoritesDB().Init()
		if err != nil {
			printErr(fmt.Errorf(theme.AppFavoritesFailedFormat, err))
			return
		}
	} else {
		af.reloadFavorList()
	}
}

func (af *appFavorites) ShortCut() fyne.Shortcut {
	return &desktop.CustomShortcut{KeyName: fyne.Key3, Modifier: fyne.KeyModifierAlt}
}

func (af *appFavorites) Icon() fyne.Resource {
	return theme.ResourceFavoritesIcon
}
