package ui

import (
	"slices"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	Task "github.com/darwinfroese/scribe/internal/task"
	"github.com/darwinfroese/scribe/internal/theme"
)

const (
	parentTreeLevel = 1
)

type tree struct {
	*tview.TreeView

	handleInput func(*tcell.EventKey) *tcell.EventKey
	focusedNode *tview.TreeNode
}

func (ui *UI) selectNextClosest(tree *tree, node *tview.TreeNode) {
	if node.GetLevel() == parentTreeLevel {
		ui.selectNextClosestParent(tree, node)
		return
	}

	ui.selectNextClosestChild(tree, node)
}

func (ui *UI) selectNextClosestParent(tree *tree, node *tview.TreeNode) {
	parents := tree.GetRoot().GetChildren()

	// we should end up with a "No Tasks!" selection when this returns
	if len(parents) == 1 {
		return
	}

	if parents[0] == node {
		tree.focusedNode = parents[1]
		tree.SetCurrentNode(parents[1])
		return
	}

	if parents[len(parents)-1] == node {
		tree.focusedNode = parents[len(parents)-2]
		tree.SetCurrentNode(parents[len(parents)-2])
		return
	}

	for idx, parent := range parents {
		if parent == node {
			tree.focusedNode = parents[idx-1]
			tree.SetCurrentNode(parents[idx-1])
			return
		}
	}
}

func (ui *UI) selectNextClosestChild(tree *tree, node *tview.TreeNode) {
	parents := tree.GetRoot().GetChildren()
	var children []*tview.TreeNode
	var parent *tview.TreeNode

	for _, prnt := range parents {
		pChildren := prnt.GetChildren()

		if slices.Contains(pChildren, node) {
			children = pChildren
			parent = prnt
			break
		}
	}

	if len(children) == 1 {
		tree.focusedNode = parent
		tree.SetCurrentNode(parent)
		return
	}

	if children[0] == node {
		tree.focusedNode = children[1]
		tree.SetCurrentNode(children[1])
		return
	}

	if children[len(children)-1] == node {
		tree.focusedNode = children[len(children)-2]
		tree.SetCurrentNode(children[len(children)-2])
		return
	}

	for idx, child := range children {
		if child == node {
			tree.focusedNode = children[idx-1]
			tree.SetCurrentNode(children[idx-1])
			return
		}
	}
}

func (ui *UI) setCurrentNode(t *tree, node *tview.TreeNode) {
	if node.GetReference() == nil {
		t.SetCurrentNode(node)
		return
	}

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
		ui.selectNextClosest(t, node)
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
	var sortOrder int

	if filter {
		sortOrder = ui.todoListSortOrder
		ids = ui.taskService.GetIncompleteTaskIDs(sortOrder)
	} else {
		sortOrder = Task.SortOrderCompletedDateDesc
		ids = ui.taskService.GetCompletedTaskIDs(sortOrder)
	}

	tree.GetRoot().ClearChildren()

	for _, id := range ids {
		if !ui.taskService.HasParent(id) {
			// children should be caught recursively
			ui.addNode(tree.GetRoot(), id, sortOrder)
		}
	}

	if len(tree.GetRoot().GetChildren()) == 0 {
		node := tview.NewTreeNode("No Tasks!").
			SetSelectedTextStyle(tcell.StyleDefault.Foreground(theme.Color(ui.theme.TextFocus)).Background(theme.Color(ui.theme.BackgroundFocus)))
		tree.GetRoot().AddChild(node)
		tree.focusedNode = node

		tree.SetCurrentNode(node)
		return
	}

	if tree.focusedNode == nil {
		tree.focusedNode = tree.GetRoot().GetChildren()[0]
	}
}

func (ui *UI) addNode(base *tview.TreeNode, id, sortOrder int) {
	text := ui.parseColors(ui.taskService.DisplayString(id))
	task := &task{id, text}

	listItemText := ui.parseColors(ui.taskService.DisplayString(id))
	node := tview.NewTreeNode(listItemText).
		SetSelectedTextStyle(tcell.StyleDefault.Foreground(theme.Color(ui.theme.TextFocus)).Background(theme.Color(ui.theme.BackgroundFocus))).
		SetReference(task)

	base.AddChild(node)

	if ui.taskService.HasChildren(id) {
		children := ui.taskService.GetChildren(id, sortOrder)

		for _, child := range children {
			ui.addNode(node, child, sortOrder)
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
	case 'S':
		ui.todoListSortOrder = Task.SortOrderPriorityDesc
		ui.refresh()

		return nil
	case 's':
		ui.todoListSortOrder = Task.SortOrderPriorityAsc
		ui.refresh()

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

		task.text = ui.parseColors(ui.taskService.DisplayString(task.id))

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
