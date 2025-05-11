package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type list struct {
	*tview.List

	handleInput func(*tcell.EventKey) *tcell.EventKey
}

func (ui *UI) refreshLists() {
	ui.todoTaskIDs = ui.refreshTaskList(ui.todoList, hideCompleted)
	ui.completedTaskIDs = ui.refreshTaskList(ui.completedList, hideIncomplete)

	ui.refreshSessionList(ui.sessionList)
}

func (ui *UI) refreshTaskList(list *list, filter bool) []int {
	originalIndex := list.GetCurrentItem()
	list.Clear()

	allIDs := ui.taskService.GetAllTaskIDs()

	newIDs := []int{}

	for _, id := range allIDs {
		if ui.taskService.IsCompleted(id) == filter {
			continue
		}

		newIDs = append(newIDs, id)
		listItemText := ui.taskService.DisplayString(id)
		list.AddItem(listItemText, "", 0, nil)
	}

	if len(newIDs) == 0 {
		list.AddItem("No Tasks!", "", 0, nil)
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

	return newIDs
}

func (ui *UI) refreshSessionList(list *list) {
	originalIndex := list.GetCurrentItem()
	list.Clear()

	sessionIDs := ui.taskService.GetAllSessionIDs()

	if len(sessionIDs) == 0 {
		list.AddItem("No Sessions!", "", 0, nil)
		return
	}

	for idx := len(sessionIDs) - 1; idx >= 0; idx-- {
		listItemText := ui.taskService.SessionDisplayString(sessionIDs[idx])
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

func (ui *UI) wipInputHandler(event *tcell.EventKey) *tcell.EventKey {
	if ui.genericListInputHandler(event) == nil {
		return nil
	}

	switch event.Rune() {
	case ' ':
		if len(ui.todoTaskIDs) == 0 {
			return nil
		}

		index := ui.todoList.GetCurrentItem()
		ui.taskService.CompleteTask(ui.todoTaskIDs[index])

		ui.refreshLists()
		return nil
	case 'e':
		index := ui.todoList.GetCurrentItem()
		if len(ui.todoTaskIDs) == 0 {
			return nil
		}

		task, priority := ui.taskService.GetTaskDetails(ui.todoTaskIDs[index])
		ui.showEditTaskForm(task, priority)
	case 'p':
		index := ui.todoList.GetCurrentItem()
		if len(ui.todoTaskIDs) == 0 {
			return nil
		}

		ui.taskService.TogglePlanTask(ui.todoTaskIDs[index])
		ui.refreshLists()
		return nil
	case 'x':
		index := ui.todoList.GetCurrentItem()
		if len(ui.todoTaskIDs) == 0 {
			return nil
		}

		ui.taskService.DeleteTask(ui.todoTaskIDs[index])
		ui.refreshLists()

		return nil

	case 'j': // down
		index := ui.todoList.GetCurrentItem()
		if index < ui.taskService.Count()-1 {
			ui.todoList.SetCurrentItem(index + 1)
		}
		return nil

	case 'k': // up
		index := ui.todoList.GetCurrentItem()
		if index > 0 {
			ui.todoList.SetCurrentItem(index - 1)
		}
		return nil
	}

	return nil
}

func (ui *UI) completeInputHandler(event *tcell.EventKey) *tcell.EventKey {
	if ui.genericListInputHandler(event) == nil {
		return nil
	}

	switch event.Rune() {
	case ' ':
		index := ui.completedList.GetCurrentItem()
		if len(ui.completedTaskIDs) == 0 {
			return nil
		}

		ui.taskService.UnCompleteTask(ui.completedTaskIDs[index])
		ui.refreshLists()

		return nil
	case 'x':
		index := ui.completedList.GetCurrentItem()
		if len(ui.completedTaskIDs) == 0 {
			return nil
		}

		ui.taskService.DeleteTask(ui.completedTaskIDs[index])
		ui.refreshLists()

		return nil

	case 'j': // down
		index := ui.completedList.GetCurrentItem()
		if index < ui.taskService.Count()-1 {
			ui.completedList.SetCurrentItem(index + 1)
		}
		return nil

	case 'k': // up
		index := ui.completedList.GetCurrentItem()
		if index > 0 {
			ui.completedList.SetCurrentItem(index - 1)
		}
		return nil
	}

	return nil
}

func (ui *UI) genericListInputHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'a':
		ui.showNewTaskForm()
		return nil

	case 'n':
		ui.showNoteForm()
		return nil

	case 'q':
		ui.app.Stop()
		return nil
	}

	return event
}
