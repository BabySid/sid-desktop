package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"sid-desktop/theme"
)

type sodorJobTaskInstance struct {
	insID int32

	tabItem *container.TabItem

	jobID           *widget.Label
	jobName         *widget.Label
	jobInstanceID   *widget.Label
	taskNames       *widget.CheckGroup
	headerContainer *fyne.Container

	instanceCard        *widget.Card
	instanceListBinding binding.UntypedList
	instanceHeader      *widget.List
	instanceList        *widget.List
}

func newSodorJobTaskInstance(id int32) *sodorJobTaskInstance {
	ins := sodorJobTaskInstance{}
	ins.insID = id

	ins.buildTaskNames()
	ins.buildTasksInstanceInfo()

	ins.tabItem = container.NewTabItem(theme.AppSodorJobInfoTaskInstance, nil)
	ins.tabItem.Content = container.NewBorder(
		ins.headerContainer, nil, nil, nil,
		ins.instanceCard)
	return &ins
}

func (s *sodorJobTaskInstance) buildTaskNames() {
	s.jobID = widget.NewLabel("12345")
	s.jobName = widget.NewLabel("jobName")
	s.jobInstanceID = widget.NewLabel(fmt.Sprintf("%d", s.insID))
	s.taskNames = widget.NewCheckGroup([]string{"task1", "task2", "task3", "task4"}, nil)
	s.taskNames.Horizontal = true

	s.headerContainer = container.NewVBox(
		container.NewHBox(
			widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoJobID, s.jobID)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoJobName, s.jobName)),
			widget.NewForm(widget.NewFormItem(theme.AppSodorJobInfoTaskInstanceJobInsID, s.jobInstanceID)),
			layout.NewSpacer()),
	)
	s.headerContainer.Add(widget.NewCard("", theme.AppSodorJobInfoTaskTitle, s.taskNames))
}

func (s *sodorJobTaskInstance) buildTasksInstanceInfo() {
	s.instanceListBinding = binding.NewUntypedList()
	s.instanceListBinding.Set([]interface{}{1, 2, 3})

	s.instanceHeader = widget.NewList(
		func() int {
			return 1
		},
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil,
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				nil,
				container.NewGridWithColumns(8,
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignTrailing, fyne.TextStyle{})),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(theme.AppSodorJobInfoTaskInstanceID)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(theme.AppSodorJobInfoTaskInstanceTaskName)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(theme.AppSodorJobInfoTaskInstanceStartTime)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText(theme.AppSodorJobInfoTaskInstanceStopTime)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*widget.Label).SetText(theme.AppSodorJobInfoTaskInstanceRunHost)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*widget.Label).SetText(theme.AppSodorJobInfoTaskInstancePID)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[5].(*widget.Label).SetText(theme.AppSodorJobInfoTaskInstanceExitCode)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[6].(*widget.Label).SetText(theme.AppSodorJobInfoTaskInstanceExitMsg)
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[7].(*widget.Label).SetText("")
		},
	)

	s.instanceList = widget.NewListWithData(
		s.instanceListBinding,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil,
				widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{}),
				nil,
				container.NewGridWithColumns(8,
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{}),
					container.NewHBox(
						layout.NewSpacer(),
						widget.NewButtonWithIcon(theme.AppSodorJobInfoTaskInstanceListOp1, theme.ResourceInstanceIcon, nil),
					),
				),
			)
		},
		func(data binding.DataItem, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText("id")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText("task name")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText("2022-12-12 20:20:33")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[2].(*widget.Label).SetText("2022-12-12 20:20:35")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[3].(*widget.Label).SetText("127.0.0.1")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[4].(*widget.Label).SetText("12345")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[5].(*widget.Label).SetText("0")
			item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[6].(*widget.Label).SetText("OK")

			op1Btn := item.(*fyne.Container).Objects[0].(*fyne.Container).Objects[7].(*fyne.Container).Objects[1].(*widget.Button)
			op1Btn.SetText(theme.AppSodorJobInfoTaskInstanceListOp1)
			op1Btn.OnTapped = func() {
				s.viewTaskInstanceDetail()
			}
		},
	)

	s.instanceCard = widget.NewCard("", theme.AppSodorJobInfoTaskInstance,
		container.NewScroll(container.NewBorder(s.instanceHeader, nil, nil, nil, s.instanceList)))
}

func (s *sodorJobTaskInstance) viewTaskInstanceDetail() {
	pContent := widget.NewMultiLineEntry()
	pContent.SetMinRowsVisible(8)
	pContent.SetText(`hello worldhdsakjhdsjdjsldjsd`)

	outVars := widget.NewMultiLineEntry()
	outVars.SetMinRowsVisible(8)
	outVars.SetText(`{"output":"12345"}`)
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
