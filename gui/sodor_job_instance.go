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
)

type sodorJobInstance struct {
	jID   int32
	jName string

	tabItem *container.TabItem

	refresh *widget.Button
	jobID   *widget.Label
	jobName *widget.Label

	instanceCard        *widget.Card
	instanceListBinding binding.UntypedList
	instanceHeader      fyne.CanvasObject
	instanceList        *widget.List

	viewTaskInstanceHandle func(taskInstanceParam)

	jobInstanceCache *sodor.JobInstances
}

type taskInstanceParam struct {
	jID    int32
	jInsID int32
}

func getTabNameOfJobInstance(info *sodor.Job) string {
	return fmt.Sprintf("%s-%d", theme.AppSodorJobInfoJobInstance, info.Id)
}

func newSodorJobInstance(job *sodor.Job) *sodorJobInstance {
	ins := sodorJobInstance{}
	ins.jID = job.Id
	ins.jName = job.Name

	ins.refresh = widget.NewButtonWithIcon(theme.AppPageRefresh, theme.ResourceRefreshIcon, func() {
		ins.loadJobInstance()
	})

	ins.buildJobInstanceInfo()

	ins.jobID = widget.NewLabel(fmt.Sprintf("%d", ins.jID))
	ins.jobName = widget.NewLabel(ins.jName)

	ins.tabItem = container.NewTabItem(getTabNameOfJobInstance(job), nil)
	ins.tabItem.Content = container.NewBorder(
		container.NewHBox(
			widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoJobID, ins.jobID)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoJobName, ins.jobName)),
			layout.NewSpacer(),
			ins.refresh,
		),
		nil, nil, nil,
		ins.instanceCard)

	go ins.loadJobInstance()
	return &ins
}

func (s *sodorJobInstance) buildJobInstanceInfo() {
	s.instanceListBinding = binding.NewUntypedList()

	s.instanceHeader = container.NewBorder(nil, nil,
		widget.NewLabelWithStyle(theme.AppSodorJobInfoJobInstanceID, fyne.TextAlignLeading, fyne.TextStyle{}),
		nil,
		container.NewGridWithColumns(6,
			widget.NewLabelWithStyle(theme.AppSodorJobInfoJobInstanceScheduleTime, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(theme.AppSodorJobInfoJobInstanceStartTime, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(theme.AppSodorJobInfoJobInstanceStopTime, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(theme.AppSodorJobInfoJobInstanceExitCode, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(theme.AppSodorJobInfoJobInstanceExitMsg, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
		),
	)

	s.instanceList = widget.NewListWithData(
		s.instanceListBinding,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil,
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				nil,
				container.NewGridWithColumns(6,
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					container.NewHBox(
						layout.NewSpacer(),
						widget.NewButtonWithIcon(theme.AppSodorJobListOp2, theme.ResourceInstanceIcon, nil),
					),
				),
			)
		},
		func(data binding.DataItem, item fyne.CanvasObject) {
			o, _ := data.(binding.Untyped).Get()
			jobIns := o.(*sodor.JobInstance)

			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%d", jobIns.Id))
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(gobase.FormatTimeStamp(int64(jobIns.ScheduleTs)))
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(gobase.FormatTimeStamp(int64(jobIns.StartTs)))
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText(gobase.FormatTimeStamp(int64(jobIns.StopTs)))
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*widget.Label).SetText(fmt.Sprintf("%d", jobIns.ExitCode))
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*widget.Label).SetText(jobIns.ExitMsg)

			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[5].(*fyne.Container).Objects[1].(*widget.Button).SetText(theme.AppSodorJobListOp2)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[5].(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
				if s.viewTaskInstanceHandle != nil {
					s.viewTaskInstanceHandle(taskInstanceParam{jID: s.jID, jInsID: jobIns.Id})
				}
			}
		},
	)

	s.instanceCard = widget.NewCard("", theme.AppSodorJobInfoJobInstance,
		container.NewScroll(container.NewBorder(s.instanceHeader, nil, nil, nil, s.instanceList)))
}

func (s *sodorJobInstance) loadJobInstance() {
	req := sodor.JobInstance{Id: s.jID}
	resp := &sodor.JobInstances{}
	err := common.GetSodorClient().Call(common.SelectJobInstances, &req, resp)
	if err != nil {
		printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
		return
	}

	s.jobInstanceCache = resp

	wrapper := common.NewJobInstanceWrapper(s.jobInstanceCache)
	s.instanceListBinding.Set(wrapper.AsInterfaceArray())
}
