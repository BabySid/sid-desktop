package gui

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

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

	taskRelation *widget.List

	okHandle      func()
	dismissHandle func()
	ok            *widget.Button
	dismiss       *widget.Button
}

func newSodorJobInfo(jID int32) *sodorJobInfo {
	info := sodorJobInfo{}
	info.jID = jID

	// job brief
	//infoBox := container.NewVBox()
	//if info.jID > 0 {
	//	info.jobID = widget.NewLabel(fmt.Sprintf("%d", info.jID))
	//}
	//
	//info.jobName = widget.NewEntry()
	//
	//infoBox.Add(container.NewHBox(
	//	container.NewHBox(widget.NewLabel(theme.AppSodorJobInfoJobID), info.jobID, widget.NewLabel(theme.AppSodorJobInfoJobName), info.jobName, layout.NewSpacer()),
	//))
	//
	//info.cronSpec = widget.NewSelect(getDefaultCronSpec(), nil)
	//info.schedulerMode = widget.NewRadioGroup(getJobScheduleMode(), func(s string) {
	//	if s == "Crontab" {
	//		info.cronSpec.Enable()
	//	} else {
	//		info.cronSpec.Disable()
	//	}
	//})
	//info.schedulerMode.Horizontal = false
	//
	//infoForm.Append(theme.AppSodorJobInfoJobScheduleMode, info.schedulerMode)
	//infoForm.Append("", info.cronSpec)
	//
	//info.alertGroup = widget.NewSelect([]string{}, nil)
	//infoForm.Append(theme.AppSodorJobInfoJobAlertGroup, info.alertGroup)
	//info.jobCard = widget.NewCard("", theme.AppSodorJobInfoTitle, infoForm)
	//
	//// tasks
	//info.addTask = widget.NewButton(theme.AppSodorJobInfoAddTask, func() {
	//
	//})
	//info.taskListBinding = binding.NewUntypedList()
	//info.taskListBinding.Set([]interface{}{1, 2, 3})
	//info.tasks = widget.NewListWithData(
	//	info.taskListBinding,
	//	func() fyne.CanvasObject {
	//		return container.NewGridWithColumns(1,
	//			widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
	//			container.NewHBox(
	//				layout.NewSpacer(),
	//				widget.NewButtonWithIcon(theme.AppSodorJobListOp1, theme.ResourceRmIcon, nil),
	//			),
	//		)
	//	},
	//	func(data binding.DataItem, item fyne.CanvasObject) {
	//		item.(*fyne.Container).Objects[0].(*widget.Label).SetText("task name")
	//
	//		item.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*widget.Button).SetText(theme.AppSodorJobListOp1)
	//		item.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
	//
	//		}
	//	},
	//)
	//
	//taskForm := widget.NewForm()
	//if info.jID > 0 {
	//	info.taskID = widget.NewLabel(fmt.Sprintf("%d", info.jID))
	//	taskForm.Append(theme.AppSodorJobInfoTaskID, info.taskID)
	//}
	//info.taskName = widget.NewEntry()
	//taskForm.Append(theme.AppSodorJobInfoTaskName, info.taskName)
	//
	//info.taskType = widget.NewSelect([]string{"Shell"}, nil)
	//taskForm.Append(theme.AppSodorJobInfoTaskType, info.taskType)
	//
	//info.runningHosts = widget.NewCheckGroup([]string{
	//	"1.2.3.4",
	//	"2.3.4.5",
	//	"3.4.5.6",
	//}, nil)
	//taskForm.Append(theme.AppSodorJobInfoRunningHost, info.runningHosts)
	//
	//info.taskContent = widget.NewMultiLineEntry()
	//info.taskContent.Wrapping = fyne.TextWrapWord
	//taskForm.Append(theme.AppSodorJobInfoTaskContent, info.taskContent)
	//
	//taskSplit := container.NewHSplit(container.NewBorder(
	//	container.NewHBox(layout.NewSpacer(), info.addTask),
	//	nil, nil, nil,
	//	info.tasks,
	//), taskForm)
	//taskSplit.SetOffset(0.3)
	//info.taskCard = widget.NewCard("", theme.AppSodorJobInfoTasksTitle, taskSplit)
	//// relations
	//
	//info.ok = widget.NewButtonWithIcon(theme.ConfirmText, fyneTheme.ConfirmIcon(), func() {
	//	if info.okHandle != nil {
	//		info.okHandle()
	//	}
	//})
	//info.dismiss = widget.NewButtonWithIcon(theme.DismissText, fyneTheme.CancelIcon(), func() {
	//	if info.dismissHandle != nil {
	//		info.dismissHandle()
	//	}
	//})
	//
	//info.tabItem = container.NewTabItem(theme.AppSodorJobInfoTitle, nil)
	//info.tabItem.Content = container.NewVBox(
	//	info.jobCard,
	//	info.taskCard,
	//	container.NewBorder(nil, nil, nil, nil,
	//		container.NewHBox(layout.NewSpacer(), info.dismiss, info.ok)))
	return &info
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
