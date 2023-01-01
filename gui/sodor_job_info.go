package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	fyneTheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"sid-desktop/theme"
)

type relationItem struct {
	from      *widget.Select
	direction *widget.Label
	to        *widget.Select

	del *widget.Button

	delHandle func(item *relationItem)
	content   *fyne.Container
}

func newRelationItem() *relationItem {
	item := relationItem{}
	item.from = widget.NewSelect([]string{
		"fromTask",
	}, nil)
	item.to = widget.NewSelect([]string{
		"toTask",
	}, nil)
	item.direction = widget.NewLabel(theme.AppSodorRelationItemDirection)

	item.del = widget.NewButton(theme.AppSodorRelationItemOp1, func() {
		if item.delHandle != nil {
			item.delHandle(&item)
		}
	})

	item.content = container.NewHBox(item.from, item.direction, item.to, layout.NewSpacer(), item.del)
	return &item
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
	runningHosts    []*widget.CheckGroup
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
	return &info
}

func (s *sodorJobInfo) buildJobBrief() {
	infoBox := container.NewVBox()
	if s.jID > 0 {
		s.jobID = widget.NewLabel(fmt.Sprintf("%d", s.jID))
	} else {
		s.jobID = widget.NewLabel("")
	}

	s.jobName = widget.NewEntry()

	s.alertGroup = widget.NewSelect([]string{}, nil)

	infoBox.Add(container.NewBorder(nil, nil,
		widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoJobID, s.jobID)), nil,
		container.NewGridWithColumns(3,
			widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoJobName, s.jobName)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoJobAlertGroup, s.alertGroup)),
			layout.NewSpacer(),
		),
	))

	s.cronSpec = widget.NewSelect(getDefaultCronSpec(), nil)
	s.schedulerMode = widget.NewRadioGroup(getJobScheduleMode(), func(opt string) {
		if opt == "Crontab" {
			s.cronSpec.Enable()
		} else {
			s.cronSpec.Disable()
		}
	})
	s.schedulerMode.SetSelected("None")
	s.schedulerMode.Horizontal = true

	infoBox.Add(container.NewHBox(widget.NewLabel(theme.AppSodorJobInfoJobScheduleMode), s.schedulerMode, s.cronSpec))

	s.jobCard = widget.NewCard("", theme.AppSodorJobInfoTitle, infoBox)
}

func (s *sodorJobInfo) buildTasks() {
	s.addTask = widget.NewButton(theme.AppSodorJobInfoAddTask, func() {

	})
	s.taskListBinding = binding.NewUntypedList()
	s.taskListBinding.Set([]interface{}{1, 2, 3})
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
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText("task name")
			item.(*fyne.Container).Objects[2].(*widget.Button).SetText(theme.AppSodorTaskListOp1)
			item.(*fyne.Container).Objects[2].(*widget.Button).OnTapped = func() {

			}
		},
	)
	s.tasks.OnSelected = func(id widget.ListItemID) {

	}

	taskBox := container.NewVBox()
	if s.jID > 0 {
		s.taskID = widget.NewLabel(fmt.Sprintf("%d", s.jID))
	} else {
		s.taskID = widget.NewLabel("")
	}
	s.taskName = widget.NewEntry()

	s.taskType = widget.NewSelect([]string{"Shell"}, nil)
	s.taskType.SetSelectedIndex(0)

	taskBox.Add(container.NewBorder(nil, nil,
		widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoTaskID, s.taskID)),
		nil,
		container.NewGridWithColumns(3,
			widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoTaskName, s.taskName)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoTaskType, s.taskType)),
			layout.NewSpacer())),
	)

	s.runningHosts = make([]*widget.CheckGroup, 0)
	for i := 0; i < 2; i++ {
		cg := widget.NewCheckGroup([]string{
			"1.2.3.4",
			"2.3.4.5",
			"3.4.5.6",
		}, nil)
		cg.Horizontal = true
		s.runningHosts = append(s.runningHosts, cg)
	}

	hosts := container.NewVBox()
	hosts.Add(s.runningHosts[0])
	hosts.Add(s.runningHosts[1])
	taskBox.Add(widget.NewAccordion(widget.NewAccordionItem(theme.AppSodorJobInfoRunningHost, hosts)))

	s.taskContent = widget.NewMultiLineEntry()
	s.taskContent.Wrapping = fyne.TextWrapWord
	s.taskContent.SetMinRowsVisible(8)
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
		item := newRelationItem()
		s.relationData = append(s.relationData, item)
		item.delHandle = func(item *relationItem) {
			s.relationList.Remove(item.content)
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
		"None",
		"Crontab",
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
