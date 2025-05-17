package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type tree struct {
	*tview.TreeView

	handleInput func(*tcell.EventKey) *tcell.EventKey
}

func (ui *UI) refreshTrees() {
	ui.todoTaskIDs = ui.refreshTaskTree(ui.todoList, hideCompleted)
	ui.completedTaskIDs = ui.refreshTaskTree(ui.completedList, hideIncomplete)
}

func (ui *UI) refreshTaskTree(tree *tree, filter bool) []*task {
	current := tree.GetCurrentNode()
	allIDs := ui.taskService.GetAllTaskIDs()

	newIDs := []*task{}
	tree.GetRoot().ClearChildren()

	for _, id := range allIDs {
		if ui.taskService.IsCompleted(id) == filter {
			continue
		}

		id := id
		text := ui.taskService.DisplayString(id)
		task := &task{id, text}

		newIDs = append(newIDs, task)
		listItemText := ui.taskService.DisplayString(id)
		node := tview.NewTreeNode(listItemText).
			SetSelectedTextStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.NewHexColor(0xffe5b3))).
			SetReference(task)

		tree.GetRoot().AddChild(node)
	}

	if len(newIDs) == 0 {
		node := tview.NewTreeNode("No Tasks!").
			SetSelectedTextStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.NewHexColor(0xffe5b3)))
		tree.GetRoot().AddChild(node)
	}

	if current == nil {
		tree.SetCurrentNode(tree.GetRoot().GetChildren()[0])
	} else {
		setCurrentNode(tree, current)
	}

	return newIDs
}

func (ui *UI) wipInputHandler(event *tcell.EventKey) *tcell.EventKey {
	if ui.genericListInputHandler(event) == nil {
		return nil
	}

	switch event.Rune() {
	case 'K': // UP
		children := ui.activeTaskList.GetRoot().GetChildren()
		selected := ui.activeTaskList.GetCurrentNode()

		for idx, child := range children {
			if child == selected {
				if idx == 0 {
					return nil
				}

				children[idx] = children[idx-1]
				children[idx-1] = selected
				ui.activeTaskList.GetRoot().SetChildren(children)

				return nil
			}
		}

		return nil
	case 'J':
		children := ui.activeTaskList.GetRoot().GetChildren()
		selected := ui.activeTaskList.GetCurrentNode()

		for idx, child := range children {
			if child == selected {
				if idx == len(children)-1 {
					return nil
				}

				children[idx] = children[idx+1]
				children[idx+1] = selected
				ui.activeTaskList.GetRoot().SetChildren(children)

				return nil
			}
		}

		return nil
	case 'e':
		task := ui.todoList.GetCurrentNode().GetReference().(*task)
		text, priority := ui.taskService.GetTaskDetails(task.id)
		ui.showEditTaskForm(text, priority)

		return nil
	case 'p':
		task := ui.todoList.GetCurrentNode().GetReference().(*task)
		ui.taskService.TogglePlanTask(task.id)

		task.text = ui.taskService.DisplayString(task.id)

		ui.todoList.GetCurrentNode().SetText(task.text)
		ui.todoList.GetCurrentNode().SetReference(task)

		return nil
	}

	return event
}

func (ui *UI) completeInputHandler(event *tcell.EventKey) *tcell.EventKey {
	return ui.genericListInputHandler(event)
}

func setCurrentNode(tree *tree, node *tview.TreeNode) {
	selectedTask := node.GetReference().(*task)

	for _, child := range tree.GetRoot().GetChildren() {
		ttask := child.GetReference().(*task)

		if ttask.id == selectedTask.id {
			tree.SetCurrentNode(child)
			return
		}
	}
}
