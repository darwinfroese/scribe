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

	SortOrderNone = iota
	SortOrderCompletedDateAsc
	SortOrderCompletedDateDesc
	SortOrderPriorityAsc
	SortOrderPriorityDesc
	SortOrderManual

	PriorityCriticalColorKey = "critical"
	PriorityHighColorKey     = "high"
	PriorityMediumColorKey   = "medium"
	PriorityLowColorKey      = "low"

	SubTextColorKey = "subtext"
)

// task is the internal task structure used for managing task details
// that shouldn't be exposed to consumers of the package. given that
// this task will be serialized and written to a file for persistence
// the fields need to be exposed.
type task struct {
	ID                int       `json:"id"`
	Completed         bool      `json:"completed"`
	Planned           bool      `json:"planned"`
	Priority          int       `json:"priority"`
	InheritedPriority int       `json:"inherited_priority"`
	Description       string    `json:"description"`
	CompletedAt       time.Time `json:"completed_at"`

	SortIndex int `json:"sort_index"`

	HasParent bool  `json:"has_parent"`
	Parent    int   `json:"parent"`
	Children  []int `json:"children"`
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

func (service *Service) AddTask(description string, priority int) {
	ttask := task{
		ID:                service.storage.Tasks.NextID,
		Description:       description,
		Priority:          priority,
		InheritedPriority: priority,
		Completed:         false,
		Planned:           false,
	}

	service.storage.Tasks.NextID++
	service.storage.Tasks.Tasks = append(service.storage.Tasks.Tasks, &ttask)

	service.write()
}

func (service *Service) AddChildTask(description string, priority int, parentDisplay string) {
	ttask := task{
		ID:                service.storage.Tasks.NextID,
		Description:       description,
		Priority:          priority,
		InheritedPriority: priority,
		Completed:         false,
		Planned:           false,
	}

	parents := service.GetAllParents()
	for _, parent := range parents {
		display := service.FormDisplayString(parent)

		if display == parentDisplay {
			parentTask := service.getTask(parent)

			ttask.Parent = parent
			ttask.HasParent = true
			parentTask.Children = append(parentTask.Children, ttask.ID)

			if ttask.Priority < parentTask.Priority && ttask.Priority < parentTask.InheritedPriority {
				parentTask.InheritedPriority = ttask.Priority
			}

			if parentTask.Planned {
				parentTask.Planned = false
				service.unplanTask(parentTask.ID)
			}

			service.updateTask(parentTask)
			break
		}
	}

	service.storage.Tasks.NextID++
	service.storage.Tasks.Tasks = append(service.storage.Tasks.Tasks, &ttask)

	service.write()
}

func (service *Service) GetAllTaskIDs() []int {
	ids := []int{}

	for _, task := range service.storage.Tasks.Tasks {
		ids = append(ids, task.ID)
	}

	return ids
}

func (service *Service) GetCompletedTaskIDs(ordering int) []int {
	ids := []int{}

	for _, task := range service.storage.Tasks.Tasks {
		if task.Completed {
			ids = append(ids, task.ID)
		}
	}

	switch ordering {
	case SortOrderCompletedDateAsc:
		slices.SortFunc(ids, service.sortOrderCompletedDateFunc(sortOrderAsc))
	case SortOrderCompletedDateDesc:
		slices.SortFunc(ids, service.sortOrderCompletedDateFunc(sortOrderDesc))
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

func (service *Service) GetIncompleteTaskIDs(ordering int) []int {
	ids := []int{}

	for _, task := range service.storage.Tasks.Tasks {
		if !task.Completed {
			ids = append(ids, task.ID)
		}
	}

	switch ordering {
	case SortOrderPriorityAsc:
		slices.SortFunc(ids, service.sortOrderPriorityFunc(sortOrderAsc))
	case SortOrderPriorityDesc:
		slices.SortFunc(ids, service.sortOrderPriorityFunc(sortOrderDesc))
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

func (service *Service) GetParent(id int) int {
	task := service.getTask(id)

	if !task.HasParent {
		return task.ID
	}

	return service.getTask(task.Parent).ID
}

func (service *Service) GetAllParents() []int {
	tasks := service.storage.Tasks.Tasks
	parents := []int{}

	for _, task := range tasks {
		if !task.HasParent {
			parents = append(parents, task.ID)
		}
	}

	return parents
}

func (service *Service) GetChildren(id, ordering int) []int {
	task := service.getTask(id)

	children := task.Children

	switch ordering {
	case SortOrderCompletedDateAsc:
		slices.SortFunc(children, service.sortOrderCompletedDateFunc(sortOrderAsc))
	case SortOrderCompletedDateDesc:
		slices.SortFunc(children, service.sortOrderCompletedDateFunc(sortOrderAsc))
	case SortOrderPriorityAsc:
		slices.SortFunc(children, service.sortOrderPriorityFunc(sortOrderAsc))
	case SortOrderPriorityDesc:
		slices.SortFunc(children, service.sortOrderPriorityFunc(sortOrderDesc))
	}

	return children
}

func (service *Service) ToggleComplete(id int) {
	for idx, task := range service.storage.Tasks.Tasks {
		if task.ID == id {
			task.Completed = !task.Completed
			task.CompletedAt = time.Now()

			if !task.Completed && len(task.Children) > 0 {
				service.resetTask(task)
				service.write()
				return
			}

			if task.HasParent {
				service.completeParent(task)
				service.adjustParentPriority(task)
			}

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

func (service *Service) AddChild(parentID, childID int) {
	parent := service.getTask(parentID)
	child := service.getTask(childID)

	// don't allow nesting more than one level
	if len(child.Children) > 0 {
		return
	}

	parent.Children = append(parent.Children, child.ID)
	child.Parent = parent.ID
	child.HasParent = true

	service.updateTask(parent)
	service.updateTask(child)

	service.write()
}

func (service *Service) RemoveChild(childID int) {
	child := service.getTask(childID)

	if !child.HasParent {
		return
	}

	parent := service.getTask(child.Parent)

	idx := slices.Index(parent.Children, child.ID)
	parent.Children = slices.Delete(parent.Children, idx, idx+1)
	child.Parent = 0
	child.HasParent = false

	if len(parent.Children) == 0 {
		parent.InheritedPriority = priorityLow
	} else {
		service.adjustParentPriority(parent)
	}

	service.updateTask(parent)
	service.updateTask(child)

	service.write()
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

	if task.HasParent {
		service.RemoveChild(task.ID)
	}

	if len(task.Children) > 0 {
		for _, childID := range task.Children {
			child := service.getTask(childID)

			child.HasParent = false
			child.Parent = 0

			service.updateTask(child)
		}
	}

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

func (service *Service) HasChildren(id int) bool {
	task := service.getTask(id)

	return len(task.Children) > 0
}

func (service *Service) HasParent(id int) bool {
	task := service.getTask(id)

	if task == nil {
		return false
	}

	return task.HasParent
}

func (service *Service) FormDisplayString(id int) string {
	task := service.getTask(id)

	return fmt.Sprintf("%d - %s", task.ID, task.Description)
}

func (service *Service) DisplayString(id int) string {
	task := service.getTask(id)

	prefix := "○"

	if task == nil {
		return "unknown task"
	}

	priority := min(task.Priority, task.InheritedPriority)
	display := fmt.Sprintf("%s [%s::](%s)[white::]", task.Description, getPriorityColor(priority), getPriorityString(priority))

	if task.Planned && service.taskPlannedToday(task.ID) {
		prefix = "→"
		display = fmt.Sprintf("[::b]%s[::B]", display)
	}

	if task.Completed {
		prefix = "✓"
		display = fmt.Sprintf("[::i]%s[::I] [%s::i]%s[white::I]", display, SubTextColorKey, task.CompletedAt.Format(time.DateOnly))
	}

	return fmt.Sprintf("%s %s", prefix, display)
}

func (service *Service) ReportString(id int) string {
	task := service.getTask(id)

	if task == nil {
		return "unknown task"
	}

	// TODO: when updating this to support parent/child, use min(task.Priority, task.InheritedPriority)
	display := fmt.Sprintf("%s (%s)", task.Description, getPriorityString(task.Priority))
	return display
}

func (service *Service) GetTaskDetails(id int) (string, int) {
	task := service.getTask(id)

	return task.Description, task.Priority
}

func (service *Service) EditTask(id int, description string, priority int) {
	task := service.getTask(id)

	task.Description = description
	task.Priority = priority

	// if we are childless or a parent
	if len(task.Children) == 0 {
		task.InheritedPriority = priority
	} else {
		task.InheritedPriority = min(task.InheritedPriority, priority)
	}

	if task.HasParent {
		parent := service.getTask(task.Parent)
		parent.InheritedPriority = min(parent.InheritedPriority, priority)
		service.updateTask(parent)
	}

	service.updateTask(task)
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

func (service *Service) adjustParentPriority(pTask *task) {
	var parent *task

	if pTask.HasParent {
		parent = service.getTask(pTask.Parent)
	} else {
		parent = pTask
	}

	if parent.Completed {
		return
	}

	highestPriority := priorityLow
	for _, child := range parent.Children {
		cTask := service.getTask(child)

		if !cTask.Completed && cTask.Priority < highestPriority {
			highestPriority = cTask.Priority
		}
	}

	parent.InheritedPriority = highestPriority

	service.updateTask(parent)
	service.write()
}

func (service *Service) completeParent(task *task) {
	parent := service.getTask(task.Parent)

	for _, child := range parent.Children {
		cTask := service.getTask(child)

		if !cTask.Completed {
			return
		}
	}

	service.ToggleComplete(parent.ID)
}

func (service *Service) getTask(id int) *task {
	for _, task := range service.storage.Tasks.Tasks {
		if task.ID == id {
			return task
		}
	}

	return nil
}

func (service *Service) updateTask(task *task) {
	for idx, tt := range service.storage.Tasks.Tasks {
		if tt.ID == task.ID {
			service.storage.Tasks.Tasks[idx] = task
			return
		}
	}
}

func (service *Service) resetTask(task *task) {
	task.Planned = false
	task.Completed = false

	for _, child := range task.Children {
		childTask := service.getTask(child)
		childTask.Planned = false
		childTask.Completed = false

		service.updateTask(childTask)
	}

	service.updateTask(task)
}

func (service *Service) write() {
	// NOTE: should this hard exit here?
	content, err := json.Marshal(service.storage)
	if err != nil {
		log.Fatal("unable to marshal the database content: ", err)
	}

	err = service.db.Write(content)
	if err != nil {
		log.Fatal("unable to write the database content: ", err)
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
		return PriorityCriticalColorKey
	case priorityHigh:
		return PriorityHighColorKey
	case priorityMedium:
		return PriorityMediumColorKey
	case priorityLow:
		return PriorityLowColorKey
	default:
		return "black"
	}
}
