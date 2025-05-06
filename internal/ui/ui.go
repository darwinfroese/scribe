package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	hideCompleted  = true
	hideIncomplete = false
)

type UI struct {
	app *tview.Application

	todoList      *list
	completedList *list
	sessionList   *list

	addTaskForm *tview.Form

	pages *tview.Pages

	taskService TaskService

	addTaskFormOpen    bool
	sessionListFocused bool

	activeTaskList *list

	todoTaskIDs      []int
	completedTaskIDs []int
}

type TaskService interface {
	AddTask(description string, priority int) int
	Count() int

	GetAllTasks() []int
	GetCompletedTasks() []int
	GetIncompleteTasks() []int

	CompleteTask(id int)
	UnCompleteTask(id int)

	IsCompleted(id int) bool
	DisplayString(id int) string
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

	ui.completedTaskIDs = ui.taskService.GetCompletedTasks()
	ui.todoTaskIDs = ui.taskService.GetIncompleteTasks()

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

	ui.todoList = &list{
		List: tview.NewList().
			ShowSecondaryText(false).
			SetSelectedFocusOnly(true).
			SetHighlightFullLine(true).
			SetSelectedStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.NewHexColor(0xffe5b3))),
		handleInput: ui.wipInputHandler,
	}
	ui.todoList.SetBorder(true).SetTitle(" Tasks (Space: Complete | a: Add | q: Quit) ")

	ui.completedList = &list{
		List: tview.NewList().
			ShowSecondaryText(false).
			SetSelectedFocusOnly(true).
			SetHighlightFullLine(true).
			SetSelectedStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.NewHexColor(0xffe5b3))),
		handleInput: ui.completeInputHandler,
	}
	ui.completedList.SetBorder(true).SetTitle(" Completed Tasks ")

	ui.pages = tview.NewPages()
	ui.addTaskForm = ui.createAddTaskForm()

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
		AddPage("form", modal(ui.addTaskForm, 100, 9), true, false)

	ui.refreshLists()

	ui.todoList.SetInputCapture(ui.listInputHandler())
	ui.completedList.SetInputCapture(ui.listInputHandler())
	ui.sessionList.SetInputCapture(ui.listInputHandler())

	ui.activeTaskList = ui.todoList
	ui.app.SetRoot(ui.pages, true)
}
