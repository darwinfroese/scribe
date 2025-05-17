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
	Completed   bool      `json:"completed"`
	Planned     bool      `json:"planned"`
	Priority    int       `json:"priority"`
	Description string    `json:"description"`
	CompletedAt time.Time `json:"completed_at"`

	Parent   int   `json:"parent"`
	Children []int `json:"children"`
}

type taskStorage struct {
	NextID int     `json:"next_id"`
	Tasks  []*task `json:"tasks"`

	DeletedTasks []*task `json:"deleted_tasks"`
}

type storage struct {
	Tasks    *taskStorage    `json:"tasks"`
	Sessions *sessionStorage `json:"sessions"`
}

type Service struct {
	db *database.Database

	storage *storage
}

func NewService(db *database.Database) *Service {
	service := Service{
		db: db,
	}

	dbContent, err := db.Read()
	if err != nil {
		log.Fatal("unable to load the database: ", err)
	}

	if len(dbContent) == 0 {
		storage := &storage{}
		storage.Tasks = &taskStorage{NextID: 0, Tasks: make([]*task, 0)}
		storage.Sessions = &sessionStorage{NextID: 0, Sessions: make([]*session, 0)}

		service.storage = storage
		return &service
	}

	storage := storage{}

	err = json.Unmarshal(dbContent, &storage)
	if err != nil {
		log.Fatal("unable to parse the database contents: ", err)
	}

	service.storage = &storage

	return &service
}

func (service *Service) AddTask(description string, priority int) int {
	ttask := task{
		ID:          service.storage.Tasks.NextID,
		Description: description,
		Priority:    priority,
		Completed:   false,
		Planned:     false,
	}

	service.storage.Tasks.NextID++
	service.storage.Tasks.Tasks = append(service.storage.Tasks.Tasks, &ttask)

	service.write()

	return ttask.ID
}

func (service *Service) GetAllTaskIDs() []int {
	ids := []int{}

	for _, task := range service.storage.Tasks.Tasks {
		ids = append(ids, task.ID)
	}

	return ids
}

func (service *Service) GetCompletedTaskIDs() []int {
	ids := []int{}

	for _, task := range service.storage.Tasks.Tasks {
		if task.Completed {
			ids = append(ids, task.ID)
		}
	}

	return ids
}

func (service *Service) GetCompletedTaskIDsForSession(id int) []int {
	ids := []int{}

	session := service.getSession(id)

	for _, task := range service.storage.Tasks.Tasks {
		if task.Completed {
			if slices.Contains(session.PlannedTasks, task.ID) {
				ids = append(ids, task.ID)
			}
		}
	}

	return ids
}

func (service *Service) GetIncompleteTaskIDs() []int {
	ids := []int{}

	for _, task := range service.storage.Tasks.Tasks {
		if !task.Completed {
			ids = append(ids, task.ID)
		}
	}

	return ids
}

func (service *Service) GetIncompleteTaskIDsForSession(id int) []int {
	ids := []int{}

	session := service.getSession(id)

	for _, task := range service.storage.Tasks.Tasks {
		if !task.Completed {
			if slices.Contains(session.PlannedTasks, task.ID) {
				ids = append(ids, task.ID)
			}
		}
	}

	return ids
}

func (service *Service) CompleteTask(id int) {
	for idx, task := range service.storage.Tasks.Tasks {
		if task.ID == id {
			task.Completed = true
			task.CompletedAt = time.Now()

			if !task.Planned {
				service.planTask(task.ID)
				task.Planned = true
			}

			service.storage.Tasks.Tasks[idx] = task
			service.write()

			return
		}
	}
}

func (service *Service) UnCompleteTask(id int) {
	for idx, task := range service.storage.Tasks.Tasks {
		if task.ID == id {
			task.Completed = false

			service.storage.Tasks.Tasks[idx] = task
			service.write()

			return
		}
	}
}

func (service *Service) DeleteTask(id int) {
	idxToDelete := 0
	taskFound := false

	for idx, task := range service.storage.Tasks.Tasks {
		if task.ID == id {
			idxToDelete = idx
			taskFound = true

			break
		}
	}

	if !taskFound {
		return
	}

	task := service.storage.Tasks.Tasks[idxToDelete]
	service.storage.Tasks.Tasks = slices.Delete(service.storage.Tasks.Tasks, idxToDelete, idxToDelete+1)
	service.storage.Tasks.DeletedTasks = append(service.storage.Tasks.DeletedTasks, task)

	service.unplanTask(task.ID)

	service.write()
}

func (service *Service) Count() int {
	return len(service.storage.Tasks.Tasks)
}

func (service *Service) IsCompleted(id int) bool {
	task := service.getTask(id)

	if task == nil {
		return false
	}

	return task.Completed
}

func (service *Service) DisplayString(id int) string {
	task := service.getTask(id)

	if task == nil {
		return "unknown task"
	}

	display := fmt.Sprintf("%s [%s::i](%s)[white::I]", task.Description, getPriorityColor(task.Priority), getPriorityString(task.Priority))

	if task.Completed {
		display = fmt.Sprintf("%s [gray::i]%s[white::I]", display, task.CompletedAt.Format(time.DateOnly))
	}

	if task.Planned && service.taskPlannedToday(task.ID) {
		display = fmt.Sprintf("[::b]%s[::B]", display)
	}

	return display
}

func (service *Service) ReportString(id int) string {
	task := service.getTask(id)

	if task == nil {
		return "unknown task"
	}

	display := fmt.Sprintf("%s (%s)", task.Description, getPriorityString(task.Priority))
	return display
}

func (service *Service) GetTaskDetails(id int) (string, int) {
	task := service.getTask(id)

	return task.Description, task.Priority
}

func (service *Service) EditTask(id int, description string, priority int) {
	for idx, task := range service.storage.Tasks.Tasks {
		if task.ID == id {
			task.Description = description
			task.Priority = priority

			service.storage.Tasks.Tasks[idx] = task
			break
		}
	}

	service.write()
}

func (service *Service) GetTasksIDsForSession(sessionID int) []int {
	for _, session := range service.storage.Sessions.Sessions {
		if session.ID == sessionID {
			return session.PlannedTasks
		}
	}

	return []int{}
}

func (service *Service) getTask(id int) *task {
	for _, task := range service.storage.Tasks.Tasks {
		if task.ID == id {
			return task
		}
	}

	return nil
}

func (service *Service) saveTask(task *task) {
	for idx, tt := range service.storage.Tasks.Tasks {
		if tt.ID == task.ID {
			service.storage.Tasks.Tasks[idx] = task
			return
		}
	}
}

func (service *Service) write() {
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
