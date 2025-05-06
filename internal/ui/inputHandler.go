package ui

import (
	"github.com/gdamore/tcell/v2"
)

func (ui *UI) listInputHandler() func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			ui.activeTaskList.handleInput(event)
		case tcell.KeyCtrlJ: // down
			if ui.addTaskFormOpen || ui.sessionListFocused {
				return event
			}

			ui.activeTaskList = ui.completedList
			ui.app.SetFocus(ui.activeTaskList)
			return nil

		case tcell.KeyCtrlK: // up
			if ui.addTaskFormOpen || ui.sessionListFocused {
				return event
			}

			ui.activeTaskList = ui.todoList
			ui.app.SetFocus(ui.activeTaskList)
			return nil

		case tcell.KeyCtrlL: // right
			if ui.addTaskFormOpen {
				return event
			}

			ui.sessionListFocused = true
			ui.app.SetFocus(ui.sessionList)
			return nil

		case tcell.KeyCtrlH: // left
			if ui.addTaskFormOpen {
				return event
			}

			ui.sessionListFocused = false
			ui.app.SetFocus(ui.activeTaskList)
			return nil

		case tcell.KeyEnter:
			return event
		}

		return event
	}
}
