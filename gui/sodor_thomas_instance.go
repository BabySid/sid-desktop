package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	"github.com/BabySid/proto/sodor"
	"sid-desktop/common"
	"sid-desktop/theme"
	sw "sid-desktop/widget"
	"strings"
	"sync"
)

type sodorThomasInstance struct {
	tid int32

	tabItem *container.TabItem

	refresh *sw.RefreshButton
	metrics *widget.Button

	thomasCard          *widget.Card
	thomasID            *widget.Label
	thomasName          *widget.Label
	thomasVersion       *widget.Label
	thomasTags          *widget.Label
	thomasProto         *widget.Label
	thomasHost          *widget.Label
	thomasPort          *widget.Label
	thomasPID           *widget.Label
	thomasStartTime     *widget.Label
	thomasHeartBeatTime *widget.Label
	thomasType          *widget.Label
	thomasStatus        *widget.Label

	instanceCard        *widget.Card
	instanceListBinding binding.UntypedList
	instanceHeader      fyne.CanvasObject
	instanceList        *widget.List

	thomasInsLock sync.Mutex
	thomasIns     *sodor.ThomasInstance
}

func newSodorThomasInstance(id int32) *sodorThomasInstance {
	ins := sodorThomasInstance{}
	ins.tid = id

	ins.refresh = sw.NewRefreshButton("*/30 * * * * *", ins.loadThomasInstance)
	ins.metrics = widget.NewButton(theme.AppPageMetrics, func() {

	})

	ins.buildThomasInfo()
	ins.buildThomasInstanceInfo()

	ins.tabItem = container.NewTabItem(theme.AppSodorThomasInfo, nil)
	ins.tabItem.Content = container.NewBorder(
		container.NewBorder(
			container.NewHBox(layout.NewSpacer(), ins.refresh.Content, ins.metrics),
			nil, nil, nil,
			ins.thomasCard),
		nil, nil, nil,
		ins.instanceCard)

	ins.loadThomasInstance()
	return &ins
}

func (s *sodorThomasInstance) buildThomasInfo() {
	s.thomasID = widget.NewLabel(fmt.Sprintf("%d", s.tid))
	s.thomasName = widget.NewLabel("")
	s.thomasVersion = widget.NewLabel("")
	s.thomasTags = widget.NewLabel("")
	s.thomasProto = widget.NewLabel("")
	s.thomasHost = widget.NewLabel("")
	s.thomasPort = widget.NewLabel("")
	s.thomasPID = widget.NewLabel("")
	s.thomasStartTime = widget.NewLabel("")
	s.thomasHeartBeatTime = widget.NewLabel("")
	s.thomasType = widget.NewLabel("")
	s.thomasStatus = widget.NewLabel("")

	infoBox := container.NewVBox()
	infoBox.Add(
		container.NewHBox(
			widget.NewForm(widget.NewFormItem(theme.AppSodorThomasInfoID, s.thomasID)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorThomasInfoName, s.thomasName)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorThomasInfoVersion, s.thomasVersion)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorThomasInfoTags, s.thomasTags)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorThomasInfoType, s.thomasType)),
		),
	)
	infoBox.Add(
		container.NewHBox(
			widget.NewForm(widget.NewFormItem(theme.AppSodorThomasInfoProto, s.thomasProto)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorThomasInfoHost, s.thomasHost)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorThomasInfoPort, s.thomasPort)),
		),
	)
	infoBox.Add(
		container.NewHBox(
			widget.NewForm(widget.NewFormItem(theme.AppSodorThomasInfoStartTime, s.thomasStartTime)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorThomasInfoHeartBeatTime, s.thomasHeartBeatTime)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorThomasInfoPID, s.thomasPID)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorThomasInfoStatus, s.thomasStatus)),
		),
	)

	s.thomasCard = widget.NewCard("", theme.AppSodorThomasInfo, infoBox)
}

func (s *sodorThomasInstance) buildThomasInstanceInfo() {
	s.instanceListBinding = binding.NewUntypedList()

	s.instanceHeader = container.NewBorder(nil, nil,
		widget.NewLabelWithStyle(theme.AppSodorThomasInfoInstanceID, fyne.TextAlignLeading, fyne.TextStyle{}),
		nil,
		container.NewBorder(nil, nil,
			widget.NewLabelWithStyle(theme.AppSodorThomasInfoInstanceCreateTime, fyne.TextAlignCenter, fyne.TextStyle{}),
			nil,
			widget.NewLabelWithStyle(theme.AppSodorThomasInfoInstanceMetrics, fyne.TextAlignCenter, fyne.TextStyle{})),
	)

	s.instanceList = widget.NewListWithData(
		s.instanceListBinding,
		func() fyne.CanvasObject {
			metrics := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{})
			metrics.Wrapping = fyne.TextWrapWord
			return container.NewBorder(nil, nil,
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				nil,
				container.NewBorder(nil, nil,
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					nil,
					metrics),
			)
		},
		func(data binding.DataItem, item fyne.CanvasObject) {
			o, _ := data.(binding.Untyped).Get()
			metrics := o.(*sodor.ThomasMetrics)

			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%d", metrics.Id))
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(gobase.FormatTimeStamp(int64(metrics.CreateAt)))

			txt, _ := metrics.Metrics.MarshalJSON()
			metricsLabel := item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label)
			metricsLabel.SetText(string(txt))

			i := common.Find(s.instanceListBinding, data)
			s.instanceList.SetItemHeight(i, metricsLabel.MinSize().Height)
		},
	)

	s.instanceCard = widget.NewCard("", theme.AppSodorThomasInfoInstance,
		container.NewBorder(
			s.instanceHeader, nil, nil, nil,
			s.instanceList))
}

func (s *sodorThomasInstance) loadThomasInstance() {
	resp := sodor.ThomasInstance{}
	req := sodor.ThomasInfo{}
	req.Id = s.tid
	err := common.GetSodorClient().Call(common.ShowThomas, &req, &resp)
	if err != nil {
		printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
		return
	}

	s.thomasInsLock.Lock()
	defer s.thomasInsLock.Unlock()
	s.thomasIns = &resp

	go s.resetGUI()
}

func (s *sodorThomasInstance) resetGUI() {
	s.thomasInsLock.Lock()
	defer s.thomasInsLock.Unlock()

	s.thomasName.SetText(s.thomasIns.Thomas.Name)
	s.thomasVersion.SetText(s.thomasIns.Thomas.Version)
	s.thomasTags.SetText(strings.Join(s.thomasIns.Thomas.Tags, common.ArraySeparator))
	s.thomasProto.SetText(s.thomasIns.Thomas.Proto)
	s.thomasHost.SetText(s.thomasIns.Thomas.Host)
	s.thomasPort.SetText(fmt.Sprintf("%d", s.thomasIns.Thomas.Port))
	s.thomasPID.SetText(fmt.Sprintf("%d", s.thomasIns.Thomas.Pid))
	s.thomasStartTime.SetText(gobase.FormatTimeStamp(int64(s.thomasIns.Thomas.StartTime)))
	s.thomasHeartBeatTime.SetText(gobase.FormatTimeStamp(int64(s.thomasIns.Thomas.HeartBeatTime)))
	s.thomasType.SetText(s.thomasIns.Thomas.ThomasType.String())
	s.thomasStatus.SetText(s.thomasIns.Thomas.Status)

	m := common.NewThomasMetricsWrapper(s.thomasIns.Metrics)
	s.instanceListBinding.Set(m.AsInterfaceArray())
}
