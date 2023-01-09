package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"sid-desktop/backend"
	"sid-desktop/storage"
	"sid-desktop/theme"
	"strings"
)

var _ sodorInterface = (*sodorFatController)(nil)

type sodorFatController struct {
	sodorAdapter

	docs *container.AppTabs

	setFatCtrl     *widget.Button
	curFatCtrlAddr *widget.Label
}

func (s *sodorFatController) CreateView() fyne.CanvasObject {
	if s.content != nil {
		return s.content
	}

	s.setFatCtrl = widget.NewButton(theme.AppSodorFatCtlSetAddr, s.setFatCtlAddr)

	ctrl, err := storage.GetAppSodorDB().LoadFatCtl()
	if err != nil {
		printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
	} else if ctrl != nil {
		if err = backend.GetSodorClient().SetFatCtrlAddr(*ctrl); err != nil {
			printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
		}
	}

	s.curFatCtrlAddr = widget.NewLabel(backend.GetSodorClient().GetFatCrl().Addr)

	s.docs = container.NewAppTabs()
	s.docs.Append(
		container.NewTabItemWithIcon(theme.AppSodorFatCtrlTabName, theme.ResourceFatCtrlIcon,
			container.NewVBox(
				container.NewBorder(nil, nil, s.setFatCtrl, nil, s.curFatCtrlAddr)),
		),
	)
	s.docs.SetTabLocation(container.TabLocationTop)

	s.content = s.docs
	return s.content
}

func (s *sodorFatController) setFatCtlAddr() {
	addr := widget.NewEntry()
	addr.SetText(backend.GetSodorClient().GetFatCrl().Addr)
	addr.Validator = validation.NewRegexp(`\S+`, theme.AppSodorFatCtlAddrValidateMsg)
	addr.SetPlaceHolder(theme.AppSodorFatCtlAddrPlaceHolder)

	cont := widget.NewForm(widget.NewFormItem(theme.AppSodorFatCtlAddr, addr))
	diag := dialog.NewCustomConfirm(theme.AppSodorFatCtlSetAddr, theme.ConfirmText, theme.DismissText, cont, func(b bool) {
		if b {
			if strings.TrimSpace(addr.Text) == "" {
				return
			}

			ctrl := backend.GetSodorClient().GetFatCrl()
			ctrl.Addr = addr.Text

			if err := backend.GetSodorClient().SetFatCtrlAddr(ctrl); err != nil {
				printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
				return
			}
			if err := storage.GetAppSodorDB().SetFatCtrl(ctrl); err != nil {
				printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
				return
			}

			s.curFatCtrlAddr.SetText(ctrl.Addr)
		}
	}, globalWin.win)

	diag.Resize(fyne.NewSize(500, 200))
	diag.Show()
}
