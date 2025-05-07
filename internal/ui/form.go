package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type formActionHandler func(form *form) func()

type form struct {
	*tview.Form

	name string
}

func (ui *UI) createForm(action string, name string, actionHandler formActionHandler) *form {
	form := &form{
		Form: tview.NewForm(),
		name: name,
	}

	taskInput := tview.NewInputField().SetLabel("Task:").SetFieldWidth(80)
	form.AddFormItem(taskInput)

	dropDown := tview.NewDropDown().SetLabel("Priority:").SetOptions([]string{"Critical", "High", "Medium", "Low"}, nil)

	dropDown.SetFocusedStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.ColorSlateGray))
	dropDown.SetListStyles(
		tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tview.Styles.PrimitiveBackgroundColor),
		tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.NewHexColor(0xffe5b3)))
	dropDown.SetPrefixStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.NewHexColor(0xffe5b3)))

	form.AddFormItem(dropDown)
	form.AddButton("Save", actionHandler(form)).
		AddButton("Cancel", func() {
			ui.hideTaskForm(name)
		})

	form.SetBorder(true).SetTitle(fmt.Sprintf("%s Task", action))

	form.SetFieldStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.ColorLightGray))
	form.SetButtonStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tview.Styles.PrimitiveBackgroundColor))
	form.SetButtonActivatedStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.ColorSlateGray))

	form.SetInputCapture(ui.formInputHandler)

	return form
}

func (ui *UI) formInputHandler(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEsc {
		ui.hideTaskForm(ui.activeForm.name)
		return nil
	}

	return event
}

func (ui *UI) showNewTaskForm() {
	ui.activeForm = ui.addTaskForm
	ui.showTaskForm("", 0, addTaskFormName)
}

func (ui *UI) showEditTaskForm(task string, priority int) {
	ui.activeForm = ui.editTaskForm
	ui.showTaskForm(task, priority, editTaskFormName)
}

func (ui *UI) showTaskForm(task string, priority int, form string) {
	taskInput := ui.activeForm.GetFormItemByLabel("Task:").(*tview.InputField)
	priorityDropDown := ui.activeForm.GetFormItemByLabel("Priority:").(*tview.DropDown)

	taskInput.SetText(task)
	priorityDropDown.SetCurrentOption(priority)

	ui.pages.ShowPage(form)
	ui.app.SetFocus(taskInput)

	ui.taskFormOpen = true
}

func (ui *UI) hideTaskForm(form string) {
	ui.pages.HidePage(form)

	if ui.sessionListFocused {
		ui.app.SetFocus(ui.sessionList)
	} else {
		ui.app.SetFocus(ui.activeTaskList)
	}

	ui.taskFormOpen = false
}

func (ui *UI) addTaskActionHandler(form *form) func() {
	return func() {
		taskDescInput := form.GetFormItemByLabel("Task:").(*tview.InputField)
		priorityDropDown := form.GetFormItemByLabel("Priority:").(*tview.DropDown)

		taskDesc := taskDescInput.GetText()
		priority, _ := priorityDropDown.GetCurrentOption()

		if taskDesc == "" {
			return
		}

		ui.todoTaskIDs = append(ui.todoTaskIDs, ui.taskService.AddTask(taskDesc, priority))
		ui.refreshLists()

		ui.hideTaskForm(addTaskFormName)
	}
}

func (ui *UI) editTaskActionHandler(form *form) func() {
	return func() {
		taskDescInput := form.GetFormItemByLabel("Task:").(*tview.InputField)
		priorityDropDown := form.GetFormItemByLabel("Priority:").(*tview.DropDown)

		taskDesc := taskDescInput.GetText()
		priority, _ := priorityDropDown.GetCurrentOption()

		if taskDesc == "" {
			return
		}

		idx := ui.todoList.GetCurrentItem()

		ui.taskService.EditTask(ui.todoTaskIDs[idx], taskDesc, priority)
		ui.refreshLists()

		ui.hideTaskForm(editTaskFormName)
	}
}
