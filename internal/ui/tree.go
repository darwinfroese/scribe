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
		// .SetSelectable()

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

	// if list.GetItemCount() > 0 {
	// 	if originalIndex >= list.GetItemCount() {
	// 		list.SetCurrentItem(list.GetItemCount() - 1)
	// 	} else if originalIndex < 0 && list.GetItemCount() > 0 {
	// 		list.SetCurrentItem(0)
	// 	} else {
	// 		list.SetCurrentItem(originalIndex)
	// 	}
	// }

	return newIDs
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

		// index := ui.todoList.GetCurrentItem()
		// ui.taskService.CompleteTask(ui.todoTaskIDs[index])
		//
		// ui.refreshLists()
		return nil
	case 'e':
		// index := ui.todoList.GetCurrentItem()
		// if len(ui.todoTaskIDs) == 0 {
		// 	return nil
		// }
		//
		// task, priority := ui.taskService.GetTaskDetails(ui.todoTaskIDs[index])
		// ui.showEditTaskForm(task, priority)
		return nil
	case 'p':
		if len(ui.todoTaskIDs) == 0 {
			return nil
		}

		task := ui.todoList.GetCurrentNode().GetReference().(*task)
		ui.taskService.TogglePlanTask(task.id)

		task.text = ui.taskService.DisplayString(task.id)

		ui.todoList.GetCurrentNode().SetText(task.text)
		ui.todoList.GetCurrentNode().SetReference(task)

		return nil
	case 'x':
		// index := ui.todoList.GetCurrentItem()
		// if len(ui.todoTaskIDs) == 0 {
		// 	return nil
		// }
		//
		// ui.taskService.DeleteTask(ui.todoTaskIDs[index])
		// ui.refreshLists()

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
		// index := ui.completedList.GetCurrentItem()
		// if len(ui.completedTaskIDs) == 0 {
		// 	return nil
		// }
		//
		// ui.taskService.UnCompleteTask(ui.completedTaskIDs[index])
		// ui.refreshLists()

		return nil
	case 'x':
		// index := ui.completedList.GetCurrentItem()
		// if len(ui.completedTaskIDs) == 0 {
		// 	return nil
		// }
		//
		// ui.taskService.DeleteTask(ui.completedTaskIDs[index])
		// ui.refreshLists()

		return nil

		// case 'j': // down
		// 	index := ui.completedList.GetCurrentItem()
		// 	if index < ui.taskService.Count()-1 {
		// 		ui.completedList.SetCurrentItem(index + 1)
		// 	}
		// 	return nil
		//
		// case 'k': // up
		// 	index := ui.completedList.GetCurrentItem()
		// 	if index > 0 {
		// 		ui.completedList.SetCurrentItem(index - 1)
		// 	}
		// 	return nil
	}

	return nil
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
