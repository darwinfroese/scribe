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
			if ui.formOpen || ui.sessionListFocused {
				return event
			}

			ui.activeTaskList = ui.completedList
			ui.focus(ui.activeTaskList)

			return nil

		case tcell.KeyCtrlK: // up
			if ui.formOpen || ui.sessionListFocused {
				return event
			}

			ui.activeTaskList = ui.todoList
			ui.focus(ui.activeTaskList)

			return nil

		case tcell.KeyCtrlL: // right
			if ui.formOpen {
				return event
			}

			ui.sessionListFocused = true
			ui.app.SetFocus(ui.sessionList)
			ui.focus(nil)
			return nil

		case tcell.KeyCtrlH: // left
			if ui.formOpen {
				return event
			}

			ui.sessionListFocused = false

			ui.focus(ui.activeTaskList)
			return nil

		case tcell.KeyEnter:
			return event
		}

		return event
	}
}

func (ui *UI) focus(tree *tree) {
	ui.todoList.SetCurrentNode(nil)
	ui.completedList.SetCurrentNode(nil)

	// so that we can clear the trees when switching
	// to session list
	if tree == nil {
		return
	}

	ui.app.SetFocus(tree)
	tree.SetCurrentNode(tree.GetRoot().GetChildren()[0])
}
