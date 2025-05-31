package ui

import (
	"slices"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	Task "github.com/darwinfroese/scribe/internal/task"
)

type tree struct {
	*tview.TreeView

	handleInput func(*tcell.EventKey) *tcell.EventKey
	focusedNode *tview.TreeNode
}

func (ui *UI) setCurrentNode(t *tree, node *tview.TreeNode) {
	task := node.GetReference().(*task)
	root := t.GetRoot()

	if ui.taskService.HasParent(task.id) {
		parentID := ui.taskService.GetParent(task.id)
		root = ui.findNode(root, parentID)
	}

	if root == nil {
		t.SetCurrentNode(t.GetRoot().GetChildren()[0])
		return
	}

	target := ui.findNode(root, task.id)

	if target != nil {
		t.SetCurrentNode(target)
	} else {
		// TODO: when we have an "order" field select the next closest in the order
		t.SetCurrentNode(root.GetChildren()[0])
	}
}

func (ui *UI) findNode(root *tview.TreeNode, id int) *tview.TreeNode {
	children := root.GetChildren()
	idx := slices.IndexFunc(children, func(node *tview.TreeNode) bool {
		nTask := node.GetReference().(*task)

		return nTask.id == id
	})

	if idx == -1 {
		return nil
	}

	return children[idx]
}

func (ui *UI) refreshTrees() {
	ui.refreshTaskTree(ui.todoList, hideCompleted)
	ui.refreshTaskTree(ui.completedList, hideIncomplete)
}

func (ui *UI) refreshTaskTree(tree *tree, filter bool) {
	var ids []int

	if filter == hideCompleted {
		ids = ui.taskService.GetIncompleteTaskIDs()
	} else {
		ids = ui.taskService.GetCompletedTaskIDs(Task.SortOrderCompletedDateDesc)
	}

	tree.GetRoot().ClearChildren()

	for _, id := range ids {
		if ui.taskService.IsCompleted(id) == filter {
			continue
		}

		if !ui.taskService.HasParent(id) {
			// children should be caught recursively
			ui.addNode(tree.GetRoot(), id)
		}
	}

	if len(tree.GetRoot().GetChildren()) == 0 {
		node := tview.NewTreeNode("No Tasks!").
			SetSelectedTextStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.NewHexColor(0xffe5b3)))
		tree.GetRoot().AddChild(node)
		tree.focusedNode = node

		tree.SetCurrentNode(node)
		return
	}

	if tree.focusedNode == nil {
		tree.focusedNode = tree.GetRoot().GetChildren()[0]
	}
}

func (ui *UI) addNode(base *tview.TreeNode, id int) {
	text := ui.taskService.DisplayString(id)
	task := &task{id, text}

	listItemText := ui.taskService.DisplayString(id)
	node := tview.NewTreeNode(listItemText).
		SetSelectedTextStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.NewHexColor(0xffe5b3))).
		SetReference(task)

	base.AddChild(node)

	if ui.taskService.HasChildren(id) {
		children := ui.taskService.GetChildren(id)

		for _, child := range children {
			ui.addNode(node, child)
		}
	}
}

func (ui *UI) wipInputHandler(event *tcell.EventKey) *tcell.EventKey {
	if ui.genericTreeInputHandler(event) == nil {
		return nil
	}

	switch event.Rune() {
	case 'A':
		selected := ui.activeTaskList.GetCurrentNode().GetReference()
		if selected == nil {
			return event
		}

		task := selected.(*task)
		parent := ui.taskService.GetParent(task.id)
		parents := ui.taskService.GetAllParents()

		ui.showNewTaskForm(true, parent, parents)
		return nil
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
		selected := ui.todoList.GetCurrentNode().GetReference()
		if selected == nil {
			return nil
		}

		task := selected.(*task)
		text, priority := ui.taskService.GetTaskDetails(task.id)
		ui.showEditTaskForm(text, priority)

		return nil
	case 'p':
		selected := ui.todoList.GetCurrentNode().GetReference()
		if selected == nil {
			return event
		}

		task := selected.(*task)
		ui.taskService.TogglePlanTask(task.id)

		task.text = ui.taskService.DisplayString(task.id)

		ui.todoList.GetCurrentNode().SetText(task.text)
		ui.todoList.GetCurrentNode().SetReference(task)

		ui.refresh()

		return nil
	case 't':
		children := ui.todoList.GetRoot().GetChildren()
		selected := ui.todoList.GetCurrentNode().GetReference()

		if selected == nil {
			return event
		}

		selectedTask := selected.(*task)

		for idx, child := range children {
			childTask := child.GetReference().(*task)
			if childTask.id == selectedTask.id {
				if idx == 0 {
					return nil
				}

				parent := children[idx-1].GetReference().(*task)

				ui.taskService.AddChild(parent.id, selectedTask.id)
				ui.refresh()

				return nil
			}
		}

		return nil
	}

	return event
}

func (ui *UI) completeInputHandler(event *tcell.EventKey) *tcell.EventKey {
	return ui.genericTreeInputHandler(event)
}
