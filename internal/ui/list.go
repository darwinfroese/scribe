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
	ui.refreshList(ui.todoList, ui.todoTaskIDs, hideCompleted)
	ui.refreshList(ui.completedList, ui.completedTaskIDs, hideIncomplete)
}

func (ui *UI) refreshList(list *list, ids []int, filter bool) {
	originalIndex := list.GetCurrentItem()
	list.Clear()

	if len(ids) == 0 {
		list.AddItem("No tasks!", "", 0, nil)
		return
	}

	for _, id := range ui.taskService.GetAllTasks() {
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

func (ui *UI) wipInputHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case ' ':
		index := ui.activeTaskList.GetCurrentItem()
		if len(ui.todoTaskIDs) == 0 {
			return nil
		}

		ui.taskService.CompleteTask(ui.todoTaskIDs[index])

		ui.completedTaskIDs = append(ui.completedTaskIDs, ui.todoTaskIDs[index])
		ui.todoTaskIDs = slices.Delete(ui.todoTaskIDs, index, index+1)

		ui.refreshLists()
		return nil
	case 'a':
		ui.showAddTaskForm()
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

	case 'q':
		ui.app.Stop()
		return nil
	}

	return nil
}

func (ui *UI) completeInputHandler(event *tcell.EventKey) *tcell.EventKey {
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
	case 'a':
		ui.showAddTaskForm()
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

	case 'q':
		ui.app.Stop()
		return nil
	}

	return nil
}
