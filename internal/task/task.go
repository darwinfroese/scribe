package task

import (
	"fmt"
	"time"
)

const (
	priorityCritical = iota
	priorityHigh
	priorityMedium
	priorityLow
)

type database struct {
	// directory is the location of the files
	directory string `json:"-"`

	Tasks []task `json:"tasks"`
}

// task is the internal task structure used for managing task details
// that shouldn't be exposed to consumers of the package. given that
// this task will be serialized and written to a file for persistence
// the fields need to be exposed.
type task struct {
	ID          int       `json:"id"`
	Parent      int       `json:"parent"`
	Completed   bool      `json:"completed"`
	Priority    int       `json:"priority"`
	Description string    `json:"description"`
	Children    []int     `json:"children"`
	CompletedAt time.Time `json:"completed_at"`
}

type TaskService struct {
	nextID int

	tasks []task
}

func NewTaskService() *TaskService {
	return &TaskService{nextID: 0, tasks: []task{}}
}

func (service *TaskService) AddTask(description string, priority int) {
	ttask := task{
		ID:          service.nextID,
		Description: description,
		Priority:    priority,
		Completed:   false,
	}

	service.nextID++

	service.tasks = append(service.tasks, ttask)
}

func (service *TaskService) GetAllTasks() []int {
	ids := []int{}

	for _, task := range service.tasks {
		ids = append(ids, task.ID)
	}

	return ids
}

func (service *TaskService) CompleteTask(id int) {
	for idx, task := range service.tasks {
		if task.ID == id {
			task.Completed = true
			task.CompletedAt = time.Now()

			service.tasks[idx] = task
			return
		}
	}
}

func (service *TaskService) Count() int {
	return len(service.tasks)
}

func (service *TaskService) IsCompleted(id int) bool {
	task := service.getTask(id)

	if task == nil {
		return false
	}

	return task.Completed
}

func (service *TaskService) DisplayString(id int) string {
	task := service.getTask(id)

	if task == nil {
		return "unknown task"
	}

	display := fmt.Sprintf("%s [%s::i](%s)[white::I]", task.Description, getPriorityColor(task.Priority), getPriorityString(task.Priority))

	if task.Completed {
		display = fmt.Sprintf("%s [gray::i]%s[white::I]", display, task.CompletedAt.Format(time.DateOnly))
	}

	return display
}

func (service *TaskService) getTask(id int) *task {
	for _, task := range service.tasks {
		if task.ID == id {
			return &task
		}
	}

	return nil
}

func getPriorityString(priority int) string {
	switch priority {
	case priorityCritical:
		return "Critical"
	case priorityHigh:
		return "High"
	case priorityMedium:
		return "Medium"
	case priorityLow:
		return "Low"
	default:
		return "Unknown"
	}
}

func getPriorityColor(priority int) string {
	switch priority {
	case priorityCritical:
		return "red"
	case priorityHigh:
		return "yellow"
	case priorityMedium:
		return "green"
	case priorityLow:
		return "blue"
	default:
		return "gray"
	}
}
