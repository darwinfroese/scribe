package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (ui *UI) createAddTaskForm() *tview.Form {
	form := tview.NewForm()

	taskInput := tview.NewInputField().SetLabel("Task:").SetFieldWidth(80)

	taskInput.SetFieldStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.ColorRed)) // Background(tcell.NewHexColor(0xffe5b3)))
	form.AddFormItem(taskInput)

	dropDown := tview.NewDropDown().SetLabel("Priority:").SetOptions([]string{"Critical", "High", "Medium", "Low"}, nil)

	dropDown.SetFocusedStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.ColorSlateGray))
	dropDown.SetListStyles(
		tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tview.Styles.PrimitiveBackgroundColor),
		tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.NewHexColor(0xffe5b3)))
	dropDown.SetPrefixStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.NewHexColor(0xffe5b3)))

	form.AddFormItem(dropDown)

	form.AddButton("Add", func() {
		taskDescInput := form.GetFormItemByLabel("Task:").(*tview.InputField)
		priorityDropDown := form.GetFormItemByLabel("Priority:").(*tview.DropDown)

		taskDesc := taskDescInput.GetText()
		priority, _ := priorityDropDown.GetCurrentOption()

		if taskDesc == "" {
			return
		}

		ui.todoTaskIDs = append(ui.todoTaskIDs, ui.taskService.AddTask(taskDesc, priority))
		ui.refreshLists()

		ui.hideAddTaskForm()
	}).
		AddButton("Cancel", func() {
			ui.hideAddTaskForm()
		})

	form.SetBorder(true).SetTitle("Add New Task")

	form.SetFieldStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.ColorLightGray))
	form.SetButtonStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tview.Styles.PrimitiveBackgroundColor))
	form.SetButtonActivatedStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.ColorSlateGray))

	form.SetInputCapture(ui.formInputHandler)

	return form
}

func (ui *UI) formInputHandler(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEsc {
		ui.hideAddTaskForm()
		return nil
	}

	return event
}

func (ui *UI) showAddTaskForm() {
	taskInput := ui.addTaskForm.GetFormItemByLabel("Task:").(*tview.InputField)
	priorityDropDown := ui.addTaskForm.GetFormItemByLabel("Priority:").(*tview.DropDown)

	taskInput.SetText("")
	priorityDropDown.SetCurrentOption(0)

	ui.pages.ShowPage("form")
	ui.app.SetFocus(taskInput)

	ui.addTaskFormOpen = true
}

func (ui *UI) hideAddTaskForm() {
	ui.pages.HidePage("form")

	if ui.sessionListFocused {
		ui.app.SetFocus(ui.sessionList)
	} else {
		ui.app.SetFocus(ui.activeTaskList)
	}

	ui.addTaskFormOpen = false
}
