package task

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/darwinfroese/scribe/internal/database"
)

const (
	priorityCritical = iota
	priorityHigh
	priorityMedium
	priorityLow
)

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

type taskStorage struct {
	NextID int    `json:"next_id"`
	Tasks  []task `json:"tasks"`

	DeletedTasks []task `json:"deleted_tasks"`
}

type TaskService struct {
	storage *taskStorage
	db      *database.Database
}

func NewTaskService(db *database.Database) *TaskService {
	service := TaskService{
		db: db,
	}

	dbContent, err := db.Read()
	if err != nil {
		log.Fatal("unable to load the database: ", err)
	}

	if len(dbContent) == 0 {
		storage := taskStorage{NextID: 0, Tasks: []task{}}
		service.storage = &storage

		return &service
	}

	storage := taskStorage{}

	err = json.Unmarshal(dbContent, &storage)
	if err != nil {
		log.Fatal("unable to parse the database contents: ", err)
	}

	service.storage = &storage
	return &service
}

func (service *TaskService) AddTask(description string, priority int) int {
	ttask := task{
		ID:          service.storage.NextID,
		Description: description,
		Priority:    priority,
		Completed:   false,
	}

	service.storage.NextID++
	service.storage.Tasks = append(service.storage.Tasks, ttask)

	service.write()

	return ttask.ID
}

func (service *TaskService) GetAllTaskIDs() []int {
	ids := []int{}

	for _, task := range service.storage.Tasks {
		ids = append(ids, task.ID)
	}

	return ids
}

func (service *TaskService) GetCompletedTaskIDs() []int {
	ids := []int{}

	for _, task := range service.storage.Tasks {
		if task.Completed {
			ids = append(ids, task.ID)
		}
	}

	return ids
}

func (service *TaskService) GetIncompleteTaskIDs() []int {
	ids := []int{}

	for _, task := range service.storage.Tasks {
		if !task.Completed {
			ids = append(ids, task.ID)
		}
	}

	return ids
}

func (service *TaskService) CompleteTask(id int) {
	for idx, task := range service.storage.Tasks {
		if task.ID == id {
			task.Completed = true
			task.CompletedAt = time.Now()

			service.storage.Tasks[idx] = task
			service.write()

			return
		}
	}
}

func (service *TaskService) UnCompleteTask(id int) {
	for idx, task := range service.storage.Tasks {
		if task.ID == id {
			task.Completed = false

			service.storage.Tasks[idx] = task
			service.write()

			return
		}
	}
}

func (service *TaskService) DeleteTask(id int) {
	idxToDelete := 0
	taskFound := false

	for idx, task := range service.storage.Tasks {
		if task.ID == id {
			idxToDelete = idx
			taskFound = true

			break
		}
	}

	if !taskFound {
		return
	}

	task := service.storage.Tasks[idxToDelete]
	service.storage.Tasks = slices.Delete(service.storage.Tasks, idxToDelete, idxToDelete+1)
	service.storage.DeletedTasks = append(service.storage.DeletedTasks, task)

	service.write()
}

func (service *TaskService) Count() int {
	return len(service.storage.Tasks)
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

func (service *TaskService) GetTaskDetails(id int) (string, int) {
	task := service.getTask(id)

	return task.Description, task.Priority
}

func (service *TaskService) EditTask(id int, description string, priority int) {
	for idx, task := range service.storage.Tasks {
		if task.ID == id {
			task.Description = description
			task.Priority = priority

			service.storage.Tasks[idx] = task
			return
		}
	}
}

func (service *TaskService) getTask(id int) *task {
	for _, task := range service.storage.Tasks {
		if task.ID == id {
			return &task
		}
	}

	return nil
}

func (service *TaskService) write() {
	// NOTE: should this hard exit here?
	content, err := json.Marshal(service.storage)
	if err != nil {
		log.Fatal("unable to marshal the database content: ", err)
	}

	err = service.db.Write(content)
	if err != nil {
		log.Fatal("unable to wirte the database content: ", err)
	}
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
