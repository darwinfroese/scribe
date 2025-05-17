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

	return nil
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
