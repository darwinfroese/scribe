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
	ui.refreshTrees()

	ui.refreshSessionList(ui.sessionList)
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
