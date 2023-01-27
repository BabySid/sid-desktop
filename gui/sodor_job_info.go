package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	fyneTheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	"github.com/BabySid/proto/sodor"
	"image/color"
	"sid-desktop/common"
	"sid-desktop/theme"
	"strings"
)

type relationItem struct {
	from      *widget.Select
	direction *widget.Label
	to        *widget.Select

	del *widget.Button

	delHandle func(item *relationItem)
	content   *fyne.Container
}

func newRelationItem(tasks []string) *relationItem {
	item := relationItem{}
	item.from = widget.NewSelect(tasks, nil)
	item.to = widget.NewSelect(tasks, nil)
	item.direction = widget.NewLabel(theme.AppSodorRelationItemDirection)

	item.del = widget.NewButton(theme.AppSodorRelationItemOp1, func() {
		if item.delHandle != nil {
			item.delHandle(&item)
		}
	})

	item.content = container.NewHBox(item.from, item.direction, item.to, layout.NewSpacer(), item.del)
	return &item
}

func (rel *relationItem) resetOption(opts []string) {
	selOpts := rel.from.SelectedIndex()
	rel.from.Options = opts
	rel.from.SetSelectedIndex(selOpts)

	selOpts = rel.to.SelectedIndex()
	rel.to.Options = opts
	rel.to.SetSelectedIndex(selOpts)
}

type sodorJobInfo struct {
	jID int32

	tabItem *container.TabItem

	jobCard       *widget.Card
	jobID         *widget.Label
	jobName       *widget.Entry
	schedulerMode *widget.RadioGroup
	cronSpec      *widget.Select
	alertGroup    *widget.Select

	taskCard        *widget.Card
	addTask         *widget.Button
	taskListBinding binding.UntypedList
	tasks           *widget.List
	taskID          *widget.Label
	taskName        *widget.Entry
	taskType        *widget.Select
	runningHosts    *widget.CheckGroup
	taskContent     *widget.Entry

	relationCard      *widget.Card
	addRelation       *widget.Button
	relationData      []*relationItem
	relationList      *fyne.Container
	relationAccordion *widget.Accordion

	okHandle      func()
	dismissHandle func()
	ok            *widget.Button
	dismiss       *widget.Button

	jobObj *sodor.Job
}

func newSodorJobInfo(jID int32) *sodorJobInfo {
	info := sodorJobInfo{}
	info.jID = jID

	// job brief
	info.buildJobBrief()

	// tasks
	info.buildTasks()

	// relations
	info.buildTaskRelations()

	// others
	info.ok = widget.NewButtonWithIcon(theme.ConfirmText, fyneTheme.ConfirmIcon(), func() {
		info.submitHandle()

		if info.okHandle != nil {
			info.okHandle()
		}
	})
	info.dismiss = widget.NewButtonWithIcon(theme.DismissText, fyneTheme.CancelIcon(), func() {
		if info.dismissHandle != nil {
			info.dismissHandle()
		}
	})

	info.tabItem = container.NewTabItem(theme.AppSodorJobInfoTitle, nil)
	info.tabItem.Content = container.NewBorder(
		nil, container.NewHBox(layout.NewSpacer(), info.dismiss, info.ok), nil, nil,
		container.NewScroll(container.NewBorder(info.jobCard, info.relationCard, nil, nil, info.taskCard)))

	go info.loadJob()

	return &info
}

func (s *sodorJobInfo) buildJobBrief() {
	infoBox := container.NewVBox()

	groups := common.GetSodorCache().GetAlertGroups()
	var groupOpt []string
	if groups != nil {
		groupOpt = make([]string, len(groups.AlertGroups))
		for i, g := range groups.AlertGroups {
			groupOpt[i] = fmt.Sprintf("%d:%s", g.Id, g.Name)
		}
	}

	s.jobName = widget.NewEntry()
	s.alertGroup = widget.NewSelect(groupOpt, nil)

	if s.jID > 0 {
		s.jobID = widget.NewLabel(fmt.Sprintf("%d", s.jID))
		infoBox.Add(container.NewBorder(nil, nil,
			widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoJobID, s.jobID)), nil,
			container.NewGridWithColumns(3,
				widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoJobName, s.jobName)),
				widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoJobAlertGroup, s.alertGroup)),
				layout.NewSpacer(),
			),
		))
	} else {
		infoBox.Add(container.NewBorder(nil, nil,
			nil, nil,
			container.NewGridWithColumns(3,
				widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoJobName, s.jobName)),
				widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoJobAlertGroup, s.alertGroup)),
				layout.NewSpacer(),
			),
		))
	}

	s.cronSpec = widget.NewSelect(getDefaultCronSpec(), nil)
	s.cronSpec.SetSelectedIndex(0)
	s.schedulerMode = widget.NewRadioGroup(getJobScheduleMode(), func(opt string) {
		if opt == sodor.ScheduleMode_ScheduleMode_Crontab.String() {
			s.cronSpec.Enable()
		} else {
			s.cronSpec.Disable()
		}
	})
	s.schedulerMode.SetSelected(sodor.ScheduleMode_ScheduleMode_Crontab.String())
	s.schedulerMode.Horizontal = true

	infoBox.Add(container.NewHBox(widget.NewLabel(theme.AppSodorJobInfoJobScheduleMode), s.schedulerMode, s.cronSpec))

	s.jobCard = widget.NewCard("", theme.AppSodorJobInfoTitle, infoBox)
}

func (s *sodorJobInfo) buildTasks() {
	s.addTask = widget.NewButton(theme.AppSodorJobInfoAddTask, func() {
		s.taskListBinding.Append(&sodor.Task{Name: fmt.Sprintf(theme.AppSodorJobInfoNewTaskFormat, gobase.FormatDateTime())})
	})

	s.taskListBinding = binding.NewUntypedList()
	s.taskListBinding.AddListener(newTaskNameListener(s))

	s.tasks = widget.NewListWithData(
		s.taskListBinding,
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				layout.NewSpacer(),
				widget.NewButtonWithIcon(theme.AppSodorTaskListOp1, theme.ResourceRmIcon, nil),
			)
		},
		func(data binding.DataItem, item fyne.CanvasObject) {
			o, _ := data.(binding.Untyped).Get()
			task := o.(*sodor.Task)

			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(task.Name)
			item.(*fyne.Container).Objects[2].(*widget.Button).SetText(theme.AppSodorTaskListOp1)
			item.(*fyne.Container).Objects[2].(*widget.Button).OnTapped = func() {
				rs, _ := s.taskListBinding.Get()
				rs = gobase.RemoveAnyFromSlice(rs, task)
				s.taskListBinding.Set(rs)
			}
		},
	)
	s.tasks.OnSelected = func(id widget.ListItemID) {
		task, _ := s.taskListBinding.GetValue(id)
		s.resetTaskUI(task.(*sodor.Task))
	}

	taskBox := container.NewVBox()
	s.taskID = widget.NewLabel("")
	s.taskName = widget.NewEntry()
	s.taskType = widget.NewSelect([]string{sodor.TaskType_TaskType_Shell.String()}, nil)
	s.taskType.SetSelectedIndex(0)

	taskBox.Add(container.NewBorder(nil, nil,
		widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoTaskID, s.taskID)),
		nil,
		container.NewGridWithColumns(3,
			widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoTaskName, s.taskName)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoTaskType, s.taskType)),
			layout.NewSpacer())),
	)

	runningOpts := make([]string, 0)
	thomasInfos := common.GetSodorCache().GetThomasInfos()
	if thomasInfos != nil {
		for _, thomas := range thomasInfos.ThomasInfos {
			runningOpts = append(runningOpts, fmt.Sprintf("%s %s", thomas.Host, strings.Join(thomas.Tags, common.ArraySeparator)))
		}
	}
	s.runningHosts = widget.NewCheckGroup(runningOpts, nil)
	s.runningHosts.Required = true
	if len(runningOpts) > 0 {
		s.runningHosts.SetSelected([]string{runningOpts[0]})
	}

	back := canvas.NewRectangle(color.Transparent)
	back.SetMinSize(s.setRunningHostContentSize(runningOpts))

	hosts := container.NewMax(back, container.NewScroll(s.runningHosts))
	taskBox.Add(widget.NewAccordion(widget.NewAccordionItem(theme.AppSodorJobInfoRunningHost, hosts)))

	s.taskContent = widget.NewMultiLineEntry()
	s.taskContent.Wrapping = fyne.TextWrapWord
	s.taskContent.SetMinRowsVisible(16)
	taskBox.Add(widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoTaskContent, s.taskContent)))

	taskSplit := container.NewHSplit(container.NewBorder(
		container.NewHBox(layout.NewSpacer(), s.addTask),
		nil, nil, nil,
		s.tasks,
	), taskBox)
	taskSplit.SetOffset(0.3)
	s.taskCard = widget.NewCard("", theme.AppSodorJobInfoTaskTitle, taskSplit)
}

func (s *sodorJobInfo) buildTaskRelations() {
	s.relationData = make([]*relationItem, 0)

	s.addRelation = widget.NewButton(theme.AppSodorJobInfoAddTaskRelation, func() {
		ts, _ := s.taskListBinding.Get()
		tasks := make([]string, len(ts))
		for i, v := range ts {
			tasks[i] = v.(*sodor.Task).Name
		}

		item := newRelationItem(tasks)
		s.relationData = append(s.relationData, item)
		item.delHandle = func(item *relationItem) {
			s.relationList.Remove(item.content)
			gobase.RemoveItemFromSlice(s.relationData, item)
		}
		s.relationList.Add(item.content)

		s.relationAccordion.OpenAll()
	})

	s.relationAccordion = widget.NewAccordion()
	s.relationList = container.NewVBox()
	s.relationAccordion.Append(widget.NewAccordionItem(theme.AppSodorJobInfoRelationTitle, s.relationList))

	s.relationCard = widget.NewCard("", "", container.NewBorder(
		container.NewHBox(layout.NewSpacer(), s.addRelation), nil, nil, nil,
		s.relationAccordion,
	))
}

func getJobScheduleMode() []string {
	return []string{
		sodor.ScheduleMode_ScheduleMode_None.String(),
		sodor.ScheduleMode_ScheduleMode_Crontab.String(),
	}
}

func getDefaultCronSpec() []string {
	return []string{
		"0 */5 * * * *",
		"0 */10 * * * *",
		"0 0 */1 * * *",
		"0 0 0 */1 * *",
	}
}

func (s *sodorJobInfo) loadJob() error {
	if s.jID == 0 {
		return nil
	}
	resp := sodor.Job{}
	req := sodor.Job{Id: s.jID}
	err := common.GetSodorClient().Call(common.SelectJob, &req, &resp)
	if err != nil {
		return err
	}
	s.jobObj = &resp

	s.resetUI()
	return nil
}

func (s *sodorJobInfo) resetUI() {
	groups := common.GetSodorCache().GetAlertGroups()
	var groupOpt []string
	if groups != nil {
		groupOpt = make([]string, len(groups.AlertGroups))
		for i, g := range groups.AlertGroups {
			groupOpt[i] = fmt.Sprintf("%d:%s", g.Id, g.Name)
		}
		s.alertGroup.Options = groupOpt
		s.alertGroup.Refresh()
	}
}

func (s *sodorJobInfo) resetTaskUI(task *sodor.Task) {
	if task.Id > 0 {
		s.taskID.SetText(fmt.Sprintf("%d", task.Id))
	}

	s.taskName.SetText(task.Name)
	s.taskContent.SetText(task.Content)
	s.taskType.SetSelected(task.Type.String())
	host := make([]string, len(task.RunningHosts))
	for i, h := range task.RunningHosts {
		host[i] = h.Node
	}

	s.runningHosts.SetSelected(host)
	if len(s.runningHosts.Selected) == 0 && len(s.runningHosts.Options) > 0 {
		s.runningHosts.SetSelected([]string{s.runningHosts.Options[0]})
	}
}

func (s *sodorJobInfo) resetTaskRelationUI(tasks []string) {
	for _, item := range s.relationData {
		item.resetOption(tasks)
	}
}

func (s *sodorJobInfo) setRunningHostContentSize(opts []string) fyne.Size {
	size := len(opts)
	if size >= 3 {
		size = 3
	}

	return fyne.NewSize(200, common.GetItemsHeightInCheck(size))
}

func (s *sodorJobInfo) submitHandle() {

}

type taskNameListener struct {
	s *sodorJobInfo
}

func newTaskNameListener(s *sodorJobInfo) *taskNameListener {
	return &taskNameListener{s: s}
}

func (t *taskNameListener) DataChanged() {
	tasks, _ := t.s.taskListBinding.Get()
	names := make([]string, len(tasks))
	for i, o := range tasks {
		task := o.(*sodor.Task)
		names[i] = task.Name
	}
	t.s.resetTaskRelationUI(names)
}
