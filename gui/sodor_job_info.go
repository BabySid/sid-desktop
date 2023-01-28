package gui

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/layout"
	fyneTheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/BabySid/gobase"
	"github.com/BabySid/proto/sodor"
	"image/color"
	"sid-desktop/common"
	"sid-desktop/theme"
	"strconv"
	"strings"
)

type relationItem struct {
	from *widget.Select
	to   *widget.Select

	del *widget.Button

	delHandle func(item *relationItem)
	content   *fyne.Container
}

func newRelationItem(tasks []string) *relationItem {
	item := relationItem{}
	item.from = widget.NewSelect(tasks, nil)
	item.from.SetSelectedIndex(0)
	item.to = widget.NewSelect(tasks, nil)
	item.to.SetSelectedIndex(0)

	item.del = widget.NewButton(theme.AppSodorRelationItemOp1, func() {
		if item.delHandle != nil {
			item.delHandle(&item)
		}
	})

	item.content = container.NewBorder(nil, nil, nil, item.del,
		container.NewGridWithColumns(3,
			widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoTaskRelationFrom, item.from)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoTaskRelationTo, item.to))),
		layout.NewSpacer())
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
	curTask         *sodor.Task
	taskBox         *fyne.Container
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

	ok      *widget.Button
	dismiss *widget.Button

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
		if err := info.submitHandle(); err == nil {
			if info.okHandle != nil {
				info.okHandle()
			}
		} else {
			printErr(fmt.Errorf(theme.ProcessSodorFailedFormat, err))
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

	s.jobName = widget.NewEntry()
	s.jobName.Validator = validation.NewRegexp(`\S+`, theme.AppSodorJobInfoJobName+" must not be empty")
	s.alertGroup = widget.NewSelect([]string{}, nil)
	s.setAlertGroupOpts()

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
		task := &sodor.Task{Name: fmt.Sprintf(theme.AppSodorJobInfoNewTaskFormat, gobase.FormatDateTime())}
		s.taskListBinding.Append(task)
		s.tasks.Select(s.taskListBinding.Length() - 1)
	})

	s.taskListBinding = binding.NewUntypedList()
	s.taskListBinding.AddListener(newTaskNameListener(s))

	s.taskBox = container.NewVBox()

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
				s.taskBox.Hide()

				if s.curTask != nil && s.curTask == task {
					s.curTask = nil
				}
			}
		},
	)

	s.tasks.OnSelected = func(id widget.ListItemID) {
		obj, _ := s.taskListBinding.GetValue(id)
		task := obj.(*sodor.Task)
		s.resetTaskUI(task)
		s.taskBox.Show()
	}

	s.taskID = widget.NewLabel("")
	s.taskName = widget.NewEntry()
	s.taskName.Validator = validation.NewRegexp(`\S+`, theme.AppSodorJobInfoTaskName+" must not be empty")
	s.taskName.OnChanged = func(_ string) {
		if s.curTask == nil {
			return
		}

		s.curTask.Name = s.taskName.Text

		for i := 0; i < s.taskListBinding.Length(); i++ {
			v, _ := s.taskListBinding.GetValue(i)
			if v.(*sodor.Task) == s.curTask {
				s.taskListBinding.SetValue(i, s.curTask)
				s.tasks.Refresh()
				break
			}
		}
	}
	s.taskType = widget.NewSelect([]string{sodor.TaskType_TaskType_Shell.String()}, nil)
	s.taskType.SetSelectedIndex(0)

	if s.jobObj != nil {
		s.taskBox.Add(container.NewBorder(nil, nil,
			widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoTaskID, s.taskID)),
			nil,
			container.NewGridWithColumns(2,
				widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoTaskName, s.taskName)),
				widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoTaskType, s.taskType)),
			)),
		)
	} else {
		s.taskBox.Add(container.NewBorder(nil, nil, nil, nil,
			container.NewGridWithColumns(2,
				widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoTaskName, s.taskName)),
				widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoTaskType, s.taskType)),
			)),
		)
	}

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
	s.taskBox.Add(widget.NewAccordion(widget.NewAccordionItem(theme.AppSodorJobInfoRunningHost, hosts)))

	s.taskContent = widget.NewMultiLineEntry()
	s.taskContent.Validator = validation.NewRegexp(`\S+`, theme.AppSodorJobInfoTaskContent+" must not be empty")
	s.taskContent.Wrapping = fyne.TextWrapWord
	s.taskContent.SetMinRowsVisible(16)
	s.taskBox.Add(widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoTaskContent, s.taskContent)))
	s.taskBox.Hide()

	taskSplit := container.NewHSplit(container.NewBorder(
		container.NewHBox(layout.NewSpacer(), s.addTask),
		nil, nil, nil,
		s.tasks,
	), s.taskBox)
	taskSplit.SetOffset(0.3)
	s.taskCard = widget.NewCard("", theme.AppSodorJobInfoTaskTitle, taskSplit)
}

func (s *sodorJobInfo) buildTaskRelations() {
	s.relationData = make([]*relationItem, 0)

	s.addRelation = widget.NewButton(theme.AppSodorJobInfoAddTaskRelation, func() {
		tasks := s.getAllTaskNames()
		if s.createRelationItem(tasks) != nil {
			s.relationAccordion.OpenAll()
		}
	})

	s.relationAccordion = widget.NewAccordion()
	s.relationList = container.NewVBox()
	s.relationAccordion.Append(widget.NewAccordionItem(theme.AppSodorJobInfoRelationTitle, s.relationList))

	s.relationCard = widget.NewCard("", "", container.NewBorder(
		container.NewHBox(layout.NewSpacer(), s.addRelation), nil, nil, nil,
		s.relationAccordion,
	))
}

func (s *sodorJobInfo) getAllTaskNames() []string {
	ts, _ := s.taskListBinding.Get()
	tasks := make([]string, len(ts))
	for i, v := range ts {
		tasks[i] = v.(*sodor.Task).Name
	}

	return tasks
}

func (s *sodorJobInfo) createRelationItem(tasks []string) *relationItem {
	if len(tasks) < 2 {
		return nil
	}

	item := newRelationItem(tasks)
	s.relationData = append(s.relationData, item)
	item.delHandle = func(item *relationItem) {
		s.relationList.Remove(item.content)
		gobase.RemoveItemFromSlice(s.relationData, item)
	}
	s.relationList.Add(item.content)
	return item
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

	s.resetUIWithJob()
	return nil
}

func (s *sodorJobInfo) resetUIWithJob() {
	// job brief
	s.jobID.SetText(fmt.Sprintf("%d", s.jobObj.Id))
	s.jobName.SetText(s.jobObj.Name)
	group := common.GetSodorCache().FindAlertGroupByID(s.jobObj.Id)
	if group != nil {
		s.alertGroup.SetSelected(fmt.Sprintf("%d:%s", group.Id, group.Name))
	}
	s.schedulerMode.SetSelected(s.jobObj.ScheduleMode.String())
	if s.jobObj.RoutineSpec != nil {
		s.cronSpec.SetSelected(s.jobObj.RoutineSpec.CtSpec)
	}

	// tasks
	taskWrapper := common.NewJobTasksWrapper(s.jobObj)
	s.taskListBinding.Set(taskWrapper.AsInterfaceArray())
	s.tasks.Select(0)

	// relation
	tasks := s.getAllTaskNames()
	for _, rel := range s.jobObj.Relations {
		item := s.createRelationItem(tasks)
		item.from.SetSelected(rel.FromTask)
		item.to.SetSelected(rel.ToTask)
	}
}

func (s *sodorJobInfo) saveCurTask() {
	s.curTask.Name = s.taskName.Text

	s.curTask.RunningHosts = make([]*sodor.Host, 0)
	hosts := s.runningHosts.Selected
	for _, h := range hosts {
		ip := gobase.SplitAndTrimSpace(h, " ")[0]
		s.curTask.RunningHosts = append(s.curTask.RunningHosts, &sodor.Host{
			Type: sodor.HostType_HostType_IP,
			Node: ip,
		})
	}
	s.curTask.Content = s.taskContent.Text
	s.curTask.Type = sodor.TaskType(sodor.TaskType_value[s.taskType.Selected])
	if s.taskID.Text != "" {
		id, _ := strconv.Atoi(s.taskID.Text)
		s.curTask.Id = int32(id)
	}
}

func (s *sodorJobInfo) resetTaskUI(task *sodor.Task) {
	if s.curTask == nil {
		s.curTask = task
	} else if s.curTask != task {
		// save last task
		s.saveCurTask()

		s.curTask = task
	}

	if task.Id > 0 {
		s.taskID.SetText(fmt.Sprintf("%d", task.Id))
	}

	s.taskName.SetText(task.Name)
	s.taskContent.SetText(task.Content)
	s.taskType.SetSelected(task.Type.String())

	host := make([]string, 0)
	for _, h := range task.RunningHosts {
		thomas := common.GetSodorCache().FindThomasInfo(h.Node)
		if thomas != nil {
			host = append(host, fmt.Sprintf("%s %s", thomas.Host, strings.Join(thomas.Tags, common.ArraySeparator)))
		}
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

func (s *sodorJobInfo) setAlertGroupOpts() {
	groups := common.GetSodorCache().GetAlertGroups()
	var groupOpt []string
	if groups != nil {
		groupOpt = make([]string, len(groups.AlertGroups))
		for i, g := range groups.AlertGroups {
			groupOpt[i] = fmt.Sprintf("%d:%s", g.Id, g.Name)
		}
	}
	s.alertGroup.Options = groupOpt
}

func (s *sodorJobInfo) submitHandle() error {
	// check job
	// check tasks
	// check relation
	s.saveCurTask()

	req := sodor.Job{}
	if s.jobObj != nil {
		req.Id = s.jobObj.Id
	}

	req.Name = s.jobName.Text
	if s.alertGroup.Selected != "" {
		idStr := gobase.SplitAndTrimSpace(s.alertGroup.Selected, ":")[0]
		id, _ := strconv.Atoi(idStr)
		group := common.GetSodorCache().FindAlertGroupByID(int32(id))
		if group == nil {
			return errors.New(fmt.Sprintf("%s", theme.AppSodorAlertGroupNotExist))
		}
		req.AlertGroupId = group.Id
	}

	if s.schedulerMode.Selected == sodor.ScheduleMode_ScheduleMode_None.String() {
		req.ScheduleMode = sodor.ScheduleMode_ScheduleMode_None
		req.RoutineSpec = nil
	} else {
		req.ScheduleMode = sodor.ScheduleMode_ScheduleMode_Crontab
		req.RoutineSpec = &sodor.RoutineSpec{CtSpec: s.cronSpec.Selected}
	}

	req.Tasks = make([]*sodor.Task, s.taskListBinding.Length())
	for i := 0; i < s.taskListBinding.Length(); i++ {
		obj, _ := s.taskListBinding.GetValue(i)
		req.Tasks[i] = obj.(*sodor.Task)
	}

	req.Relations = make([]*sodor.TaskRelation, len(s.relationData))
	for i, rel := range s.relationData {
		req.Relations[i] = &sodor.TaskRelation{
			FromTask: rel.from.Selected,
			ToTask:   rel.to.Selected,
		}
	}

	return nil
	// resp := sodor.JobReply{}
	//if s.jobObj != nil {
	//	return common.GetSodorClient().Call(common.UpdateJob, &req, &resp)
	//}
	//return common.GetSodorClient().Call(common.CreateJob, &req, &resp)
}

type taskNameListener struct {
	s *sodorJobInfo
}

func newTaskNameListener(s *sodorJobInfo) *taskNameListener {
	return &taskNameListener{s: s}
}

func (t *taskNameListener) DataChanged() {
	names := t.s.getAllTaskNames()
	if len(names) < 2 {
		// clear relation
		if t.s.relationList != nil {
			t.s.relationList.RemoveAll()
		}
		if t.s.relationData != nil {
			t.s.relationData = make([]*relationItem, 0)
		}
		return
	}

	t.s.resetTaskRelationUI(names)
}
