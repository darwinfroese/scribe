package internal

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	hideCompleted  = true
	hideIncomplete = false
)

func Run() {
	app := tview.NewApplication()

	list := tview.NewList().
		ShowSecondaryText(false).
		SetSelectedFocusOnly(true)
	list.SetBorder(true).SetTitle(" Tasks (Space: Complete | a: Add | q: Quit) ")

	completedList := tview.NewList().
		ShowSecondaryText(false).
		SetSelectedFocusOnly(true)
	completedList.SetBorder(true).SetTitle(" Completed Tasks ")

	pages := tview.NewPages()
	addTaskForm := createAddTaskForm(app, list, pages)

	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewGrid().
			SetColumns(0, width, 0).
			SetRows(0, height, 0).
			AddItem(p, 1, 1, 1, 1, 0, 0, true)
	}

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(list, 0, 3, true).
		AddItem(completedList, 0, 1, true)

	pages.
		AddPage("list", flex, true, true).
		AddPage("form", modal(addTaskForm, 50, 12), true, false)

	refreshTaskListUI(list, hideCompleted)
	refreshTaskListUI(completedList, hideIncomplete)

	list.SetInputCapture(listInputHandler(app, list, completedList, addTaskForm, pages))

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(fmt.Sprintf("Error running application: %v", err))
	}
}

// refreshTaskListUI clears and repopulates the tview.List widget
// based on the current state of the global 'tasks' slice.
func refreshTaskListUI(list *tview.List, filter bool) {
	originalIndex := list.GetCurrentItem()
	list.Clear()

	if len(taskList) == 0 {
		list.AddItem("No tasks!", "", 0, nil)
		return
	}

	for _, task := range taskList {
		if task.Completed == filter {
			continue
		}

		priority := getPriorityString(task.Priority)
		listItemText := fmt.Sprintf("%s [yellow::i](%s)[white::I]", task.Description, priority)
		list.AddItem(listItemText, "", 0, nil)
	}

	if list.GetItemCount() > 0 {
		if originalIndex >= list.GetItemCount() {
			list.SetCurrentItem(list.GetItemCount() - 1)
		} else if originalIndex < 0 && list.GetItemCount() > 0 {
			list.SetCurrentItem(0)
		} else {
			list.SetCurrentItem(originalIndex)
		}
	}
}

// listInputHandler handles key presses when the taskList has focus.
func listInputHandler(app *tview.Application, todoList *tview.List, completedList *tview.List, form *tview.Form, pages *tview.Pages) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case ' ':
				index := todoList.GetCurrentItem()
				if index >= 0 && index < len(taskList) {
					taskList[index].Completed = true
					refreshTaskListUI(todoList, hideCompleted)
					refreshTaskListUI(completedList, hideIncomplete)
				}
				return nil

			case 'a':
				showAddTaskForm(app, form, pages)
				return nil

			case 'j':
				// down
				index := todoList.GetCurrentItem()
				if index < len(taskList)-1 {
					todoList.SetCurrentItem(index + 1)
				}

			case 'k':
				// up
				index := todoList.GetCurrentItem()
				if index > 0 {
					todoList.SetCurrentItem(index - 1)
				}

			case 'q':
				app.Stop()
				return nil
			}

		case tcell.KeyEnter:
			return event
		}

		return event
	}
}

func createAddTaskForm(app *tview.Application, list *tview.List, pages *tview.Pages) *tview.Form {
	form := tview.NewForm()

	form.
		AddInputField("Task:", "", 40, nil, nil).
		AddDropDown("Priority:", []string{"Critical", "High", "Medium", "Low"}, 0, nil).
		AddButton("Add", func() {
			taskDescInput := form.GetFormItemByLabel("Task:").(*tview.InputField)
			priorityDropDown := form.GetFormItemByLabel("Priority:").(*tview.DropDown)

			taskDesc := taskDescInput.GetText()
			priority, _ := priorityDropDown.GetCurrentOption()

			if taskDesc == "" {
				return
			}

			newTask := task{Description: taskDesc, Priority: priority, Completed: false}
			taskList = append(taskList, newTask)

			refreshTaskListUI(list, hideCompleted)

			hideAddTaskForm(app, list, pages)
		}).
		AddButton("Cancel", func() {
			hideAddTaskForm(app, list, pages)
		})

	form.SetBorder(true).SetTitle("Add New Task")
	return form
}

// showAddTaskForm makes the add task form visible and focuses it.
func showAddTaskForm(app *tview.Application, form *tview.Form, pages *tview.Pages) {
	taskInput := form.GetFormItemByLabel("Task:").(*tview.InputField)
	priorityDropDown := form.GetFormItemByLabel("Priority:").(*tview.DropDown)

	taskInput.SetText("")
	priorityDropDown.SetCurrentOption(0)

	pages.ShowPage("form")
	app.SetFocus(taskInput)
}

// hideAddTaskForm hides the add task form and returns focus to the task list.
func hideAddTaskForm(app *tview.Application, list *tview.List, pages *tview.Pages) {
	pages.HidePage("form")
	app.SetFocus(list)
}
