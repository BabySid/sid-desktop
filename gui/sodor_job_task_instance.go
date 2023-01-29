package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	"github.com/BabySid/proto/sodor"
	"image/color"
	"sid-desktop/common"
	"sid-desktop/theme"
	"strconv"
)

type sodorJobTaskInstance struct {
	jobInfo taskInstanceParam

	tabItem *container.TabItem

	headerContainer *fyne.Container
	refresh         *widget.Button
	jobID           *widget.Label
	jobName         *widget.Label
	jobInstanceID   *widget.Label
	allInstance     *widget.Check
	taskNames       *widget.CheckGroup

	instanceCard        *widget.Card
	instanceListBinding binding.UntypedList
	instanceHeader      fyne.CanvasObject
	instanceList        *widget.List

	jobObj     *sodor.Job
	taskInsObj *sodor.TaskInstances

	taskIdCache map[int32]*sodor.Task
}

func newSodorJobTaskInstance(param taskInstanceParam) *sodorJobTaskInstance {
	ins := sodorJobTaskInstance{}
	ins.jobInfo = param

	ins.refresh = widget.NewButton(theme.AppPageRefresh, func() {
		ins.reload()
	})

	ins.buildTaskNames()
	ins.buildTasksInstanceInfo()

	ins.tabItem = container.NewTabItem(theme.AppSodorJobInfoTaskInstance, nil)
	ins.tabItem.Content = container.NewBorder(
		ins.headerContainer, nil, nil, nil,
		ins.instanceCard)

	go ins.reload()
	return &ins
}

func (s *sodorJobTaskInstance) buildTaskNames() {
	s.jobID = widget.NewLabel(fmt.Sprintf("%d", s.jobInfo.jID))
	s.jobName = widget.NewLabel("")
	s.jobInstanceID = widget.NewLabel(fmt.Sprintf("%d", s.jobInfo.jInsID))
	s.allInstance = widget.NewCheck(theme.AppSodorJobInfoAllTaskInstance, func(b bool) {
		s.loadTaskInstance()
	})
	s.taskNames = widget.NewCheckGroup([]string{}, func(strings []string) {
		s.searchInstanceByTask(strings)
	})
	s.taskNames.Horizontal = true

	s.headerContainer = container.NewVBox(
		container.NewHBox(
			widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoJobID, s.jobID)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoJobName, s.jobName)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoTaskInstanceJobInsID, s.jobInstanceID)),
			s.allInstance,
			layout.NewSpacer(),
			s.refresh),
	)
	s.headerContainer.Add(widget.NewCard("", theme.AppSodorJobInfoTaskTitle, s.taskNames))
}

func (s *sodorJobTaskInstance) buildTasksInstanceInfo() {
	s.instanceListBinding = binding.NewUntypedList()

	opContainer := widget.NewButtonWithIcon(theme.AppSodorJobInfoTaskInstanceListOp1, theme.ResourceInstanceIcon, nil)
	spaceLabel := canvas.NewRectangle(color.Transparent)
	size := opContainer.MinSize()
	spaceLabel.SetMinSize(fyne.NewSize(size.Width, size.Height))

	s.instanceHeader = container.NewBorder(nil, nil,
		widget.NewLabelWithStyle(theme.AppSodorJobInfoTaskInstanceID, fyne.TextAlignLeading, fyne.TextStyle{}),
		spaceLabel,
		container.NewGridWithColumns(8,
			widget.NewLabelWithStyle(theme.AppSodorJobInfoTaskInstanceJobInsID, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(theme.AppSodorJobInfoTaskInstanceTaskName, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(theme.AppSodorJobInfoTaskInstanceStartTime, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(theme.AppSodorJobInfoTaskInstanceStopTime, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(theme.AppSodorJobInfoTaskInstanceRunHost, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(theme.AppSodorJobInfoTaskInstancePID, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(theme.AppSodorJobInfoTaskInstanceExitCode, fyne.TextAlignCenter, fyne.TextStyle{}),
			widget.NewLabelWithStyle(theme.AppSodorJobInfoTaskInstanceExitMsg, fyne.TextAlignCenter, fyne.TextStyle{}),
		),
	)

	s.instanceList = widget.NewListWithData(
		s.instanceListBinding,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil,
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				opContainer,
				container.NewGridWithColumns(8,
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
				),
			)
		},
		func(data binding.DataItem, item fyne.CanvasObject) {
			o, _ := data.(binding.Untyped).Get()
			taskIns := o.(*sodor.TaskInstance)

			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%d", taskIns.Id))
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(fmt.Sprintf("%d", taskIns.JobInstanceId))
			if taskInfo := s.taskIdCache[taskIns.TaskId]; taskInfo != nil {
				item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(taskInfo.Name)
			}
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText(gobase.FormatTimeStamp(int64(taskIns.StartTs)))
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*widget.Label).SetText(gobase.FormatTimeStamp(int64(taskIns.StopTs)))
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*widget.Label).SetText(taskIns.Host)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[5].(*widget.Label).SetText(fmt.Sprintf("%d", taskIns.Pid))
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[6].(*widget.Label).SetText(fmt.Sprintf("%d", taskIns.ExitCode))
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[7].(*widget.Label).SetText(taskIns.ExitMsg)

			op1Btn := item.(*fyne.Container).Objects[2].(*widget.Button)
			op1Btn.SetText(theme.AppSodorJobInfoTaskInstanceListOp1)
			op1Btn.OnTapped = func() {
				s.viewTaskInstanceDetail(taskIns)
			}
		},
	)

	s.instanceCard = widget.NewCard("", theme.AppSodorJobInfoTaskInstance,
		container.NewBorder(s.instanceHeader, nil, nil, nil, s.instanceList))
}

func (s *sodorJobTaskInstance) viewTaskInstanceDetail(ins *sodor.TaskInstance) {
	pContent := widget.NewMultiLineEntry()
	pContent.SetMinRowsVisible(8)
	pContent.SetText(ins.ParsedContent)

	outVars := widget.NewMultiLineEntry()
	outVars.SetMinRowsVisible(8)
	bs, _ := ins.OutputVars.MarshalJSON()
	outVars.SetText(string(bs))

	parsedContent := widget.NewFormItem(theme.AppSodorJobInfoTaskInstanceParsedContent, pContent)
	outputVars := widget.NewFormItem(theme.AppSodorJobInfoTaskInstanceOutputVars, outVars)
	cont := container.NewBorder(
		nil, nil, nil, nil,
		widget.NewForm(parsedContent, outputVars),
	)
	diag := dialog.NewCustom(theme.AppSodorJobInfoTaskInstance, theme.ConfirmText, cont, globalWin.win)
	diag.Resize(fyne.NewSize(700, 300))
	diag.Show()
}

func (s *sodorJobTaskInstance) reload() {
	s.loadJobInfo()
	s.loadTaskInstance()
}

func (s *sodorJobTaskInstance) loadTaskInstance() {
	resp := sodor.TaskInstances{}
	insID := int32(0)
	if !s.allInstance.Checked {
		insID = s.jobInfo.jInsID
	}
	req := sodor.TaskInstance{JobId: s.jobInfo.jID, JobInstanceId: insID}
	err := common.GetSodorClient().Call(common.SelectTaskInstances, &req, &resp)
	if err != nil {
		printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
		return
	}

	s.taskInsObj = &resp

	s.searchInstanceByTask(s.taskNames.Selected)
}

func (s *sodorJobTaskInstance) loadJobInfo() {
	resp := sodor.Job{}
	req := sodor.Job{Id: s.jobInfo.jID}
	err := common.GetSodorClient().Call(common.SelectJob, &req, &resp)
	if err != nil {
		printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
		return
	}

	s.jobObj = &resp
	s.taskIdCache = make(map[int32]*sodor.Task)
	for _, t := range s.jobObj.Tasks {
		s.taskIdCache[t.Id] = t
	}

	s.resetUI()
}

func (s *sodorJobTaskInstance) resetUI() {
	s.jobName.SetText(s.jobObj.Name)
	tasks := make([]string, len(s.jobObj.Tasks))
	for i, t := range s.jobObj.Tasks {
		tasks[i] = fmt.Sprintf("%d:%s", t.Id, t.Name)
	}
	s.taskNames.Options = tasks
	if len(s.taskNames.Selected) == 0 {
		s.taskNames.SetSelected(tasks)
	} else {
		s.taskNames.Refresh()
	}
}

func (s *sodorJobTaskInstance) searchInstanceByTask(names []string) {
	if s.taskInsObj == nil {
		return
	}
	ids := make([]int32, len(names))
	for _, t := range names {
		id, _ := strconv.Atoi(gobase.SplitAndTrimSpace(t, ":")[0])
		ids = append(ids, int32(id))
	}
	wrapper := common.NewTaskInstanceWrapper(s.taskInsObj)
	s.instanceListBinding.Set(wrapper.AsInterfaceArray(ids...))
}
