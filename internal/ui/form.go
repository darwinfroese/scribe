package ui

import (
	"fmt"

	"github.com/darwinfroese/scribe/internal/theme"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type formActionHandler func(form *form) func()

type form struct {
	*tview.Form

	name string
}

func (ui *UI) createForm(action string, name string, childForm bool, actionHandler formActionHandler) *form {
	form := &form{
		Form: tview.NewForm(),
		name: name,
	}

	if childForm {
		parent := tview.NewDropDown().SetLabel("Parent:").SetOptions([]string{}, nil)

		parent.SetFocusedStyle(
			tcell.StyleDefault.Foreground(theme.Color(ui.theme.TextFocus)).
				Background(theme.Color(ui.theme.BackgroundFocus)),
		)

		parent.SetListStyles(
			// unselected
			tcell.StyleDefault.Foreground(theme.Color(ui.theme.Text)).Background(theme.Color(ui.theme.Background)),
			// selected
			tcell.StyleDefault.Foreground(theme.Color(ui.theme.TextFocus)).Background(theme.Color(ui.theme.BackgroundFocus)))

		parent.SetPrefixStyle(tcell.StyleDefault.Foreground(theme.Color(ui.theme.Text)).Background(theme.Color(ui.theme.Background)))

		form.AddFormItem(parent)
	}

	taskInput := tview.NewInputField().SetLabel("Task:").SetFieldWidth(80)
	form.AddFormItem(taskInput)

	dropDown := tview.NewDropDown().SetLabel("Priority:").SetOptions([]string{"Critical", "High", "Medium", "Low"}, nil)

	dropDown.SetFocusedStyle(tcell.StyleDefault.Foreground(theme.Color(ui.theme.TextFocus)).Background(theme.Color(ui.theme.BackgroundFocus)))
	dropDown.SetListStyles(
		// unselected
		tcell.StyleDefault.Foreground(theme.Color(ui.theme.Text)).Background(theme.Color(ui.theme.Background)),
		// selected
		tcell.StyleDefault.Foreground(theme.Color(ui.theme.TextFocus)).Background(theme.Color(ui.theme.BackgroundFocus)))
	dropDown.SetPrefixStyle(tcell.StyleDefault.Foreground(theme.Color(ui.theme.Text)).Background(theme.Color(ui.theme.Background)))

	form.AddFormItem(dropDown)
	form.AddButton("Save", actionHandler(form)).
		AddButton("Cancel", func() {
			ui.hideForm(name)
		})

	form.SetBorder(true).SetTitle(fmt.Sprintf(" %s Task ", action))

	form.SetFieldStyle(tcell.StyleDefault.Foreground(theme.Color(ui.theme.Text)).Background(theme.Color(ui.theme.InputBackground)))
	form.SetButtonStyle(tcell.StyleDefault.Foreground(theme.Color(ui.theme.Text)).Background(theme.Color(ui.theme.Background)))
	form.SetButtonActivatedStyle(tcell.StyleDefault.Foreground(theme.Color(ui.theme.TextFocus)).Background(theme.Color(ui.theme.BackgroundFocus)))

	form.SetInputCapture(ui.formInputHandler)

	return form
}

func (ui *UI) createNoteForm(name string, actionHandler formActionHandler) *form {
	form := &form{
		Form: tview.NewForm(),
		name: name,
	}

	// TODO: style the text area
	textArea := tview.NewTextArea()
	form.AddFormItem(textArea)

	form.AddButton("Save", actionHandler(form)).
		AddButton("Cancel", func() {
			ui.hideForm(name)
		})

	form.SetBorder(true).SetTitle(" Notes ")

	form.SetFieldStyle(tcell.StyleDefault.Foreground(theme.Color(ui.theme.Text)).Background(theme.Color(ui.theme.InputBackground)))
	form.SetButtonStyle(tcell.StyleDefault.Foreground(theme.Color(ui.theme.Text)).Background(theme.Color(ui.theme.Background)))
	form.SetButtonActivatedStyle(tcell.StyleDefault.Foreground(theme.Color(ui.theme.TextFocus)).Background(theme.Color(ui.theme.BackgroundFocus)))

	form.SetInputCapture(ui.formInputHandler)

	return form
}

func (ui *UI) formInputHandler(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEsc {
		ui.hideForm(ui.activeForm.name)
		return nil
	}

	return event
}

func (ui *UI) showNewTaskForm(child bool, parentID int, parents []int) {
	if !child {
		ui.activeForm = ui.addTaskForm
		ui.showTaskForm("", 0, addTaskFormName)

		return
	}

	parentDrop := ui.addChildTaskForm.GetFormItemByLabel("Parent:").(*tview.DropDown)

	options := []string{}
	selected := 0

	for idx, parent := range parents {
		if parent == parentID {
			selected = idx
		}

		options = append(options, ui.taskService.FormDisplayString(parent))
	}

	parentDrop.SetOptions(options, nil)
	parentDrop.SetCurrentOption(selected)

	ui.activeForm = ui.addChildTaskForm
	ui.showTaskForm("", 0, addChildTaskFormName)
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

	ui.formOpen = true
}

func (ui *UI) hideForm(form string) {
	ui.pages.HidePage(form)

	if ui.sessionListFocused {
		ui.app.SetFocus(ui.sessionList)
	} else {
		ui.app.SetFocus(ui.activeTaskList)
	}

	ui.formOpen = false
}

func (ui *UI) showNoteForm() {
	ui.activeForm = ui.addNoteForm

	note := ui.taskService.GetNote()
	input := ui.activeForm.GetFormItem(0).(*tview.TextArea)
	input.SetText(note, true)
	ui.app.SetFocus(ui.activeForm)
	ui.activeForm.SetFocus(0)

	ui.pages.ShowPage(noteFormName)

	ui.formOpen = true
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

		ui.taskService.AddTask(taskDesc, priority)
		ui.refresh()

		ui.hideForm(addTaskFormName)
	}
}

func (ui *UI) addChildTaskActionHandler(form *form) func() {
	return func() {
		parentDropDown := form.GetFormItemByLabel("Parent:").(*tview.DropDown)
		taskDescInput := form.GetFormItemByLabel("Task:").(*tview.InputField)
		priorityDropDown := form.GetFormItemByLabel("Priority:").(*tview.DropDown)

		taskDesc := taskDescInput.GetText()
		priority, _ := priorityDropDown.GetCurrentOption()
		_, parent := parentDropDown.GetCurrentOption()

		if taskDesc == "" {
			return
		}

		ui.taskService.AddChildTask(taskDesc, priority, parent)
		ui.refresh()

		ui.hideForm(addChildTaskFormName)
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

		task := ui.todoList.GetCurrentNode().GetReference().(*task)

		ui.taskService.EditTask(task.id, taskDesc, priority)
		ui.refresh()

		ui.hideForm(editTaskFormName)
	}
}

func (ui *UI) addNoteActionHandler(form *form) func() {
	return func() {
		input := form.GetFormItem(0).(*tview.TextArea)

		contents := input.GetText()

		if contents == "" {
			return
		}

		ui.taskService.SaveNote(contents)
		ui.refresh()

		ui.hideForm(noteFormName)
	}
}
