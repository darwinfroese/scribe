package ui

import (
	"github.com/gdamore/tcell/v2"
)

// this is our entrypoint for handling all inputs
func (ui *UI) listInputHandler() func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		if !ui.sessionListFocused {
			// if sessionListFocused is true then CurrentNode will be nil
			ui.focusCurrentNode()
		}

		// let the child lists active inputHandler be called first
		if ui.activeTaskList.handleInput(event) == nil {
			return nil
		}

		// TODO: ui.sessionList.handleInput(event)

		switch event.Key() {
		case tcell.KeyCtrlJ: // down
			if ui.formOpen || ui.sessionListFocused {
				return event
			}

			ui.focusCurrentNode()
			ui.activeTaskList = ui.completedList
			ui.focus(ui.activeTaskList)

			return nil

		case tcell.KeyCtrlK: // up
			if ui.formOpen || ui.sessionListFocused {
				return event
			}

			ui.focusCurrentNode()
			ui.activeTaskList = ui.todoList
			ui.focus(ui.activeTaskList)

			return nil

		case tcell.KeyCtrlL: // right
			if ui.formOpen {
				return event
			}

			ui.focusCurrentNode()
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

func (ui *UI) focusCurrentNode() {
	ui.activeTaskList.focusedNode = ui.activeTaskList.GetCurrentNode()
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
	ui.setCurrentNode(tree, tree.focusedNode)
}
