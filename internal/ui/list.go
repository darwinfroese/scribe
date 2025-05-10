package ui

import (
	"slices"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type list struct {
	*tview.List

	handleInput func(*tcell.EventKey) *tcell.EventKey
}

func (ui *UI) refreshLists() {
	ui.refreshTaskList(ui.todoList, ui.todoTaskIDs, hideCompleted)
	ui.refreshTaskList(ui.completedList, ui.completedTaskIDs, hideIncomplete)

	ui.refreshSessionList(ui.sessionList, ui.sessionIDs)
}

func (ui *UI) refreshTaskList(list *list, ids []int, filter bool) {
	originalIndex := list.GetCurrentItem()
	list.Clear()

	if len(ids) == 0 {
		list.AddItem("No tasks!", "", 0, nil)
		return
	}

	for _, id := range ui.taskService.GetAllTaskIDs() {
		if ui.taskService.IsCompleted(id) == filter {
			continue
		}

		listItemText := ui.taskService.DisplayString(id)
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

func (ui *UI) refreshSessionList(list *list, ids []int) {
	originalIndex := list.GetCurrentItem()
	list.Clear()

	if len(ids) == 0 {
		list.AddItem("No Sessions!", "", 0, nil)
		return
	}

	for _, id := range ui.taskService.GetAllSessionIDs() {
		listItemText := ui.taskService.SessionDisplayString(id)
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
		index := ui.todoList.GetCurrentItem()
		if len(ui.todoTaskIDs) == 0 {
			return nil
		}

		ui.taskService.CompleteTask(ui.todoTaskIDs[index])

		ui.completedTaskIDs = append(ui.completedTaskIDs, ui.todoTaskIDs[index])
		ui.todoTaskIDs = slices.Delete(ui.todoTaskIDs, index, index+1)

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
		ui.todoTaskIDs = slices.Delete(ui.todoTaskIDs, index, index+1)
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

		ui.todoTaskIDs = append(ui.todoTaskIDs, ui.completedTaskIDs[index])
		ui.completedTaskIDs = slices.Delete(ui.completedTaskIDs, index, index+1)

		ui.refreshLists()

		return nil
	case 'x':
		index := ui.completedList.GetCurrentItem()
		if len(ui.completedTaskIDs) == 0 {
			return nil
		}

		ui.taskService.DeleteTask(ui.completedTaskIDs[index])
		ui.completedTaskIDs = slices.Delete(ui.completedTaskIDs, index, index+1)
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
