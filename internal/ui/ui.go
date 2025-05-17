package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	hideCompleted  = true
	hideIncomplete = false

	addTaskFormName  = "add-form"
	editTaskFormName = "edit-form"

	noteFormName = "notes"
)

type UI struct {
	app *tview.Application

	todoList      *tree
	completedList *tree
	sessionList   *list

	addTaskForm  *form
	editTaskForm *form
	addNoteForm  *form

	pages *tview.Pages

	taskService TaskService

	formOpen           bool
	sessionListFocused bool

	activeTaskList *tree
	activeForm     *form

	sessionIDs []int
}

type task struct {
	id   int
	text string
}

type TaskService interface {
	AddTask(description string, priority int)
	Count() int

	GetAllTaskIDs() []int
	GetCompletedTaskIDs() []int
	GetIncompleteTaskIDs() []int
	GetTaskDetails(id int) (string, int)

	GetChildren(id int) []int

	ToggleComplete(id int)
	AddChild(parentID int, childID int)
	DeleteTask(id int)
	EditTask(id int, description string, priority int)

	TogglePlanTask(id int)

	IsCompleted(id int) bool
	HasChildren(id int) bool
	HasParent(id int) bool

	DisplayString(id int) string

	GetAllSessionIDs() []int
	SessionDisplayString(id int) string

	SaveNote(contents string)
	GetNote() string
}

func New(taskService TaskService) *UI {
	theme := tview.Styles

	theme.PrimitiveBackgroundColor = tcell.NewHexColor(0xfff0d1)
	theme.BorderColor = tcell.NewHexColor(0x0065ad)
	theme.PrimaryTextColor = tcell.NewHexColor(0x1a0b00)
	theme.SecondaryTextColor = tcell.NewHexColor(0x1a0b00)
	theme.TertiaryTextColor = tcell.NewHexColor(0x1a0b00)
	theme.TitleColor = tcell.NewHexColor(0x1a0b00)
	tview.Styles = theme

	ui := &UI{
		taskService: taskService,
	}

	ui.sessionIDs = ui.taskService.GetAllSessionIDs()
	ui.loadTasks()

	ui.build()

	return ui
}

func (ui *UI) Run() {
	if err := ui.app.Run(); err != nil {
		panic(fmt.Sprintf("Error running application: %v", err))
	}
}

func (ui *UI) build() {
	ui.app = tview.NewApplication()

	ui.todoList = &tree{
		TreeView:    createTree(),
		handleInput: ui.wipInputHandler,
	}
	ui.todoList.SetBorder(true).SetTitle(" Todo Tasks ")

	ui.completedList = &tree{
		TreeView:    createTree(),
		handleInput: ui.completeInputHandler,
	}
	ui.completedList.SetBorder(true).SetTitle(" Completed Tasks ")

	ui.pages = tview.NewPages()
	ui.addTaskForm = ui.createForm("Add New", addTaskFormName, ui.addTaskActionHandler)
	ui.editTaskForm = ui.createForm("Edit", editTaskFormName, ui.editTaskActionHandler)

	ui.addNoteForm = ui.createNoteForm(noteFormName, ui.addNoteActionHandler)

	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewGrid().
			SetColumns(0, width, 0).
			SetRows(0, height, 0).
			AddItem(p, 1, 1, 1, 1, 0, 0, true)
	}

	ui.sessionList = &list{
		List: tview.NewList().
			ShowSecondaryText(false).
			SetSelectedFocusOnly(true).
			SetHighlightFullLine(true).
			SetSelectedStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.NewHexColor(0xffe5b3))),
	}
	ui.sessionList.SetBorder(true).SetTitle(" Sessions ")

	taskFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(ui.todoList, 0, 3, true).
		AddItem(ui.completedList, 0, 1, true)

	flex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(taskFlex, 0, 3, true).
		AddItem(ui.sessionList, 0, 1, true)

	ui.pages.
		AddPage("list", flex, true, true).
		AddPage(addTaskFormName, modal(ui.addTaskForm, 100, 9), true, false).
		AddPage(editTaskFormName, modal(ui.editTaskForm, 100, 9), true, false).
		AddPage(noteFormName, modal(ui.addNoteForm, 100, 11), true, false)

	ui.activeTaskList = ui.todoList
	ui.refresh()

	ui.todoList.SetInputCapture(ui.listInputHandler())
	ui.completedList.SetInputCapture(ui.listInputHandler())
	ui.sessionList.SetInputCapture(ui.listInputHandler())

	ui.activeTaskList = ui.todoList
	ui.app.SetRoot(ui.pages, true)
}

func (ui *UI) loadTasks() {
	completedTaskIDs := ui.taskService.GetCompletedTaskIDs()
	todoTaskIDs := ui.taskService.GetIncompleteTaskIDs()

	ui.createTasksFromIDs(completedTaskIDs)
	ui.createTasksFromIDs(todoTaskIDs)
}

func (ui *UI) createTasksFromIDs(ids []int) []*task {
	tasks := []*task{}

	for _, id := range ids {
		text := ui.taskService.DisplayString(id)

		tasks = append(tasks, &task{id, text})
	}

	return tasks
}

func createTree() *tview.TreeView {
	root := tview.NewTreeNode("").
		SetSelectable(false)

	baseTree := tview.NewTreeView().
		SetGraphics(false).
		SetPrefixes([]string{"", "- ", "  - "}).
		SetRoot(root)

	return baseTree
}
