package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"sid-desktop/theme"
)

var _ sodorInterface = (*sodorFatController)(nil)

type sodorFatController struct {
	sodorAdapter

	docs *container.AppTabs

	addrBinding    binding.String
	setFatCtrl     *widget.Button
	curFatCtrlAddr *widget.Label
}

func (s *sodorFatController) CreateView() fyne.CanvasObject {
	if s.content != nil {
		return s.content
	}

	s.setFatCtrl = widget.NewButton(theme.AppSodorFatCtlSetAddr, s.setFatCtlAddr)

	s.addrBinding = binding.NewString()
	s.curFatCtrlAddr = widget.NewLabelWithData(s.addrBinding)

	s.docs = container.NewAppTabs()
	s.docs.Append(
		container.NewTabItem(theme.AppSodorFatCtrlTabName,
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
	addr.SetPlaceHolder(theme.AppSodorFatCtlAddrPlaceHolder)

	cont := widget.NewForm(widget.NewFormItem(theme.AppSodorFatCtlAddr, addr))
	diag := dialog.NewCustomConfirm(theme.AppSodorFatCtlSetAddr, theme.ConfirmText, theme.DismissText, cont, func(b bool) {
		if b {
			s.addrBinding.Set(addr.Text)
		}
	}, globalWin.win)

	diag.Resize(fyne.NewSize(500, 200))
	diag.Show()
}
