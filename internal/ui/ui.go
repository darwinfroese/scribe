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

	todoList      *tview.List
	completedList *tview.List
	sessionList   *tview.List

	addTaskForm *tview.Form

	pages *tview.Pages

	taskService TaskService
}

type TaskService interface {
	AddTask(description string, priority int)
	GetAllTasks() []int
	CompleteTask(id int)
	Count() int

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

	ui.todoList = tview.NewList().
		ShowSecondaryText(false).
		SetSelectedFocusOnly(true).
		SetHighlightFullLine(true).
		SetSelectedStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.NewHexColor(0xffe5b3)))
	ui.todoList.SetBorder(true).SetTitle(" Tasks (Space: Complete | a: Add | q: Quit) ")

	ui.completedList = tview.NewList().
		ShowSecondaryText(false).
		SetSelectedFocusOnly(true).
		SetHighlightFullLine(true).
		SetSelectedStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.NewHexColor(0xffe5b3)))
	ui.completedList.SetBorder(true).SetTitle(" Completed Tasks ")

	ui.pages = tview.NewPages()
	ui.addTaskForm = ui.createAddTaskForm()

	modal := func(p tview.Primitive, width, height int) tview.Primitive {
		return tview.NewGrid().
			SetColumns(0, width, 0).
			SetRows(0, height, 0).
			AddItem(p, 1, 1, 1, 1, 0, 0, true)
	}

	ui.sessionList = tview.NewList().
		ShowSecondaryText(false).
		SetSelectedFocusOnly(true).
		SetHighlightFullLine(true).
		SetSelectedStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.NewHexColor(0xffe5b3)))
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

	ui.refreshTaskListUI(ui.todoList, hideCompleted)
	ui.refreshTaskListUI(ui.completedList, hideIncomplete)

	ui.todoList.SetInputCapture(ui.listInputHandler())

	ui.app.SetRoot(ui.pages, true)
}

// refreshTaskListUI clears and repopulates the tview.List widget
// based on the current state of the global 'tasks' slice.
func (ui *UI) refreshTaskListUI(list *tview.List, filter bool) {
	originalIndex := list.GetCurrentItem()
	list.Clear()

	if ui.taskService.Count() == 0 {
		list.AddItem("No tasks!", "", 0, nil)
		return
	}

	for _, id := range ui.taskService.GetAllTasks() {
		if ui.taskService.IsCompleted(id) == filter {
			continue
		}

		listItemText := ui.taskService.DisplayString(id)
		list.AddItem(listItemText, "", 0, nil)
	}

	if list.GetItemCount() > 0 {
		if originalIndex >= list.GetItemCount() {
			list.SetCurrentItem(list.GetItemCount() - 1)
		} else if originalIndex < 0 && list.GetItemCount() > 0 {
			list.SetCurrentItem(0)
		} else {
			list.SetCurrentItem(originalIndex)
		}
	}
}

// listInputHandler handles key presses when the taskList has focus.
func (ui *UI) listInputHandler() func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case ' ':
				index := ui.todoList.GetCurrentItem()
				if index >= 0 && index < ui.taskService.Count() {
					ui.taskService.CompleteTask(index)
					ui.refreshTaskListUI(ui.todoList, hideCompleted)
					ui.refreshTaskListUI(ui.completedList, hideIncomplete)
				}
				return nil

			case 'a':
				ui.showAddTaskForm()
				return nil

			case 'j':
				// down
				index := ui.todoList.GetCurrentItem()
				if index < ui.taskService.Count()-1 {
					ui.todoList.SetCurrentItem(index + 1)
				}

			case 'k':
				// up
				index := ui.todoList.GetCurrentItem()
				if index > 0 {
					ui.todoList.SetCurrentItem(index - 1)
				}

			case 'q':
				ui.app.Stop()
				return nil
			}

		case tcell.KeyEnter:
			return event
		}

		return event
	}
}

func (ui *UI) createAddTaskForm() *tview.Form {
	form := tview.NewForm()

	taskInput := tview.NewInputField().SetLabel("Task:").SetFieldWidth(80)

	taskInput.SetFieldStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.ColorRed)) // Background(tcell.NewHexColor(0xffe5b3)))
	form.AddFormItem(taskInput)

	dropDown := tview.NewDropDown().SetLabel("Priority:").SetOptions([]string{"Critical", "High", "Medium", "Low"}, nil)

	dropDown.SetFocusedStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.ColorSlateGray))
	dropDown.SetListStyles(
		tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tview.Styles.PrimitiveBackgroundColor),
		tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.NewHexColor(0xffe5b3)))
	dropDown.SetPrefixStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.NewHexColor(0xffe5b3)))

	form.AddFormItem(dropDown)

	form.AddButton("Add", func() {
		taskDescInput := form.GetFormItemByLabel("Task:").(*tview.InputField)
		priorityDropDown := form.GetFormItemByLabel("Priority:").(*tview.DropDown)

		taskDesc := taskDescInput.GetText()
		priority, _ := priorityDropDown.GetCurrentOption()

		if taskDesc == "" {
			return
		}

		ui.taskService.AddTask(taskDesc, priority)
		ui.refreshTaskListUI(ui.todoList, hideCompleted)

		ui.hideAddTaskForm()
	}).
		AddButton("Cancel", func() {
			ui.hideAddTaskForm()
		})

	form.SetBorder(true).SetTitle("Add New Task")

	form.SetFieldStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.ColorLightGray))
	form.SetButtonStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tview.Styles.PrimitiveBackgroundColor))
	form.SetButtonActivatedStyle(tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor).Background(tcell.ColorSlateGray))

	form.SetInputCapture(ui.formInputHandler)

	return form
}

func (ui *UI) formInputHandler(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEsc {
		ui.hideAddTaskForm()
		return nil
	}

	return event
}

func (ui *UI) showAddTaskForm() {
	taskInput := ui.addTaskForm.GetFormItemByLabel("Task:").(*tview.InputField)
	priorityDropDown := ui.addTaskForm.GetFormItemByLabel("Priority:").(*tview.DropDown)

	taskInput.SetText("")
	priorityDropDown.SetCurrentOption(0)

	ui.pages.ShowPage("form")
	ui.app.SetFocus(taskInput)
}

func (ui *UI) hideAddTaskForm() {
	ui.pages.HidePage("form")
	ui.app.SetFocus(ui.todoList)
}
