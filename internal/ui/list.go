package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type list struct {
	*tview.List

	handleInput func(*tcell.EventKey) *tcell.EventKey
}

func (ui *UI) refresh() {
	ui.refreshSessionList(ui.sessionList)
	ui.refreshTrees()

	ui.focus(ui.activeTaskList)
}

func (ui *UI) refreshSessionList(list *list) {
	originalIndex := list.GetCurrentItem()
	list.Clear()

	sessionIDs := ui.taskService.GetAllSessionIDs(true)

	if len(sessionIDs) == 0 {
		list.AddItem("No Sessions!", "", 0, nil)
		return
	}

	for _, id := range sessionIDs {
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

func (ui *UI) genericTreeInputHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case ' ':
		selected := ui.activeTaskList.GetCurrentNode().GetReference()
		if selected == nil {
			return event
		}

		task := selected.(*task)
		ui.taskService.ToggleComplete(task.id)

		ui.refresh()

		return nil
	case 'x':
		selected := ui.activeTaskList.GetCurrentNode().GetReference()
		if selected == nil {
			return event
		}

		task := selected.(*task)
		ui.taskService.DeleteTask(task.id)

		ui.refresh()

		return nil
	case 'a':
		ui.showNewTaskForm(false, 0, nil)
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
