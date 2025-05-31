package task

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

type session struct {
	ID           int    `json:"id"`
	Date         string `json:"date"`
	Note         string `json:"note"`
	PlannedTasks []int  `json:"planned_tasks"`
}

type sessionStorage struct {
	NextID   int        `json:"next_id"`
	Sessions []*session `json:"sessions"`
}

func (service *Service) GetAllSessionIDs(reverse bool) []int {
	ids := []int{}

	for _, session := range service.storage.Sessions.Sessions {
		ids = append(ids, session.ID)
	}

	if reverse {
		slices.Reverse(ids)
	}

	return ids
}

func (service *Service) GetAllSessionDates() []string {
	dates := []string{}

	for _, session := range service.storage.Sessions.Sessions {
		dates = append(dates, session.Date)
	}

	return dates
}

func (service *Service) GetSessionDate(id int) string {
	session := service.getSession(id)

	return session.Date
}

func (service *Service) TogglePlanTask(taskID int) {
	task := service.getTask(taskID)

	task.Planned = !task.Planned
	service.saveTask(task)

	if task.Planned {
		service.planTask(taskID)
	} else {
		service.unplanTask(taskID)
	}

	if task.HasParent {
		service.planParent(task.Parent)
	}

	if len(task.Children) > 0 {
		service.planAllChildren(task.Planned, task.ID, task.Children)
	}

	service.write()
}

func (service *Service) SessionDisplayString(id int) string {
	session := service.getSession(id)

	format := ""

	if !session.isToday() {
		format = "i"
	} else {
		format = "b"
	}

	return fmt.Sprintf("[::%s]%s[::%s] ",
		format,
		service.SessionDisplayStringPlainText(id),
		strings.ToUpper(format),
	)
}

func (service *Service) SessionDisplayStringPlainText(id int) string {
	session := service.getSession(id)

	completedTasks := len(service.GetCompletedTaskIDsForSession(session.ID))

	return fmt.Sprintf("%s (%d/%d)",
		session.Date,
		completedTasks,
		len(session.PlannedTasks),
	)
}

func (service *Service) SaveNote(contents string) {
	session := service.getOrCreateTodaysSession()

	session.Note = contents

	service.saveSession(session)
	service.write()
}

func (service *Service) GetNote() string {
	session := service.getOrCreateTodaysSession()

	return session.Note
}

func (service *Service) GetNoteForSession(id int) string {
	session := service.getSession(id)

	return session.Note
}

func (service *Service) planTask(taskID int) {
	session := service.getOrCreateTodaysSession()

	if slices.Contains(session.PlannedTasks, taskID) {
		return
	}

	session.PlannedTasks = append(session.PlannedTasks, taskID)
	service.saveSession(session)
}

func (service *Service) unplanTask(taskID int) {
	session := service.getOrCreateTodaysSession()

	for idx, task := range session.PlannedTasks {
		if task == taskID {
			session.PlannedTasks = slices.Delete(session.PlannedTasks, idx, idx+1)
			service.saveSession(session)
			return
		}
	}
}

func (service *Service) saveSession(newSession *session) {
	for idx, session := range service.storage.Sessions.Sessions {
		if session.ID == newSession.ID {
			service.storage.Sessions.Sessions[idx] = newSession
			return
		}
	}
}

func (service *Service) getSession(id int) *session {
	for _, session := range service.storage.Sessions.Sessions {
		if session.ID == id {
			return session
		}
	}

	return nil
}

func (service *Service) getOrCreateTodaysSession() *session {
	today := time.Now().Format(time.DateOnly)

	for _, session := range service.storage.Sessions.Sessions {
		if session.Date == today {
			return session
		}
	}

	session := session{
		ID:   service.storage.Sessions.NextID,
		Date: today,
	}

	service.storage.Sessions.NextID++
	service.storage.Sessions.Sessions = append(service.storage.Sessions.Sessions, &session)

	return &session
}

func (service *Service) taskPlannedToday(id int) bool {
	session := service.getOrCreateTodaysSession()

	if len(session.PlannedTasks) == 0 {
		return false
	}

	if slices.Contains(session.PlannedTasks, id) {
		return true
	}

	return false
}

func (service *Service) planParent(parentID int) {
	parent := service.getTask(parentID)

	parentPlanned := true
	for _, childID := range parent.Children {
		child := service.getTask(childID)
		if !child.Planned {
			parentPlanned = false
			break
		}
	}

	parent.Planned = parentPlanned
	if parentPlanned {
		service.planTask(parent.ID)
	} else {
		service.unplanTask(parent.ID)
	}
	service.saveTask(parent)
}

func (service *Service) planAllChildren(planned bool, parentID int, children []int) {
	parent := service.getTask(parentID)
	parent.Planned = planned

	if parent.Planned {
		service.planTask(parent.ID)
	} else {
		service.unplanTask(parent.ID)
	}

	for _, childID := range children {
		child := service.getTask(childID)
		if child.Planned == parent.Planned {
			continue
		}

		child.Planned = planned

		if child.Planned {
			service.planTask(childID)
		} else {
			service.unplanTask(childID)
		}

		service.saveTask(child)
	}

	service.saveTask(parent)
}

func (session *session) isToday() bool {
	return session.Date == time.Now().Format(time.DateOnly)
}
