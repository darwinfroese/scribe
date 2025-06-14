package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	Task "github.com/darwinfroese/scribe/internal/task"
	"github.com/darwinfroese/scribe/internal/theme"
)

const (
	hideCompleted  = true
	hideIncomplete = false

	addTaskFormName      = "add-form"
	addChildTaskFormName = "add-child-form"
	editTaskFormName     = "edit-form"

	noteFormName = "notes"
)

type UI struct {
	app *tview.Application

	todoList      *tree
	completedList *tree
	sessionList   *list

	addTaskForm      *form
	addChildTaskForm *form
	editTaskForm     *form
	addNoteForm      *form

	pages *tview.Pages

	taskService TaskService

	formOpen           bool
	sessionListFocused bool

	todoListSortOrder int

	activeTaskList *tree
	activeForm     *form

	sessionIDs []int

	theme *theme.Theme
}

type task struct {
	id   int
	text string
}

type TaskService interface {
	AddTask(description string, priority int)
	AddChildTask(description string, priority int, parentDisplay string)
	Count() int

	GetAllTaskIDs() []int
	GetCompletedTaskIDs(sortOrder int) []int
	GetIncompleteTaskIDs(sortOrder int) []int
	GetTaskDetails(id int) (string, int)

	GetParent(id int) int
	GetAllParents() []int
	GetChildren(id, sortOrder int) []int

	ToggleComplete(id int)
	AddChild(parentID int, childID int)
	RemoveChild(childID int)
	DeleteTask(id int)
	EditTask(id int, description string, priority int)

	TogglePlanTask(id int)

	IsCompleted(id int) bool
	HasChildren(id int) bool
	HasParent(id int) bool

	FormDisplayString(id int) string
	DisplayString(id int) string

	GetAllSessionIDs(reverse bool) []int
	SessionDisplayString(id int) string

	SaveNote(contents string)
	GetNote() string
}

func New(taskService TaskService, userTheme *theme.Theme) *UI {
	style := tview.Styles

	style.PrimitiveBackgroundColor = theme.Color(userTheme.Background)

	style.BorderColor = theme.Color(userTheme.Border)

	style.PrimaryTextColor = theme.Color(userTheme.Text)
	style.SecondaryTextColor = theme.Color(userTheme.Text)
	style.TertiaryTextColor = theme.Color(userTheme.Text)
	style.TitleColor = theme.Color(userTheme.Text)

	tview.Styles = style

	ui := &UI{
		taskService:       taskService,
		todoListSortOrder: Task.SortOrderNone,
		theme:             userTheme,
	}

	ui.sessionIDs = ui.taskService.GetAllSessionIDs(true)
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
	ui.addTaskForm = ui.createForm("Add New", addTaskFormName, false, ui.addTaskActionHandler)
	ui.addChildTaskForm = ui.createForm("Add New Child", addChildTaskFormName, true, ui.addChildTaskActionHandler)
	ui.editTaskForm = ui.createForm("Edit", editTaskFormName, false, ui.editTaskActionHandler)

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
			SetSelectedStyle(
				tcell.StyleDefault.
					Foreground(theme.Color(ui.theme.TextFocus)).
					Background(theme.Color(ui.theme.BackgroundFocus))),
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
		AddPage(addChildTaskFormName, modal(ui.addChildTaskForm, 100, 11), true, false).
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
	completedTaskIDs := ui.taskService.GetCompletedTaskIDs(Task.SortOrderCompletedDateDesc)
	todoTaskIDs := ui.taskService.GetIncompleteTaskIDs(ui.todoListSortOrder)

	ui.createTasksFromIDs(completedTaskIDs)
	ui.createTasksFromIDs(todoTaskIDs)
}

func (ui *UI) createTasksFromIDs(ids []int) []*task {
	tasks := []*task{}

	for _, id := range ids {
		text := ui.parseColors(ui.taskService.DisplayString(id))

		tasks = append(tasks, &task{id, text})
	}

	return tasks
}

func (ui *UI) parseColors(text string) string {
	text = strings.ReplaceAll(text, fmt.Sprintf("%s::", Task.PriorityCriticalColorKey), fmt.Sprintf("%s::", ui.theme.PriorityCritical))
	text = strings.ReplaceAll(text, fmt.Sprintf("%s::", Task.PriorityHighColorKey), fmt.Sprintf("%s::", ui.theme.PriorityHigh))
	text = strings.ReplaceAll(text, fmt.Sprintf("%s::", Task.PriorityMediumColorKey), fmt.Sprintf("%s::", ui.theme.PriorityMedium))
	text = strings.ReplaceAll(text, fmt.Sprintf("%s::", Task.PriorityLowColorKey), fmt.Sprintf("%s::", ui.theme.PriorityLow))
	text = strings.ReplaceAll(text, fmt.Sprintf("%s::", Task.SubTextColorKey), fmt.Sprintf("%s::", ui.theme.SubText))

	return text
}

func createTree() *tview.TreeView {
	root := tview.NewTreeNode("").
		SetSelectable(false)

	baseTree := tview.NewTreeView().
		SetGraphics(false).
		SetPrefixes([]string{""}). // we handle the prefixes in the DisplayString func
		SetRoot(root)

	return baseTree
}
