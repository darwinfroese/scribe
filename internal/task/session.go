package task

import (
	"fmt"
	"slices"
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

func (service *Service) GetAllSessionIDs() []int {
	ids := []int{}

	for _, session := range service.storage.Sessions.Sessions {
		ids = append(ids, session.ID)
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

	service.write()
}

func (service *Service) SessionDisplayString(id int) string {
	session := service.getSession(id)

	if session.isToday() && len(session.PlannedTasks) == 0 && session.Note == "" {
		return fmt.Sprintf("[::i][::b]%s[::B][::I]", session.Date)
	}

	if session.isToday() && len(session.PlannedTasks) > 0 {
		return fmt.Sprintf("* %s", session.Date)
	}

	return session.Date
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

func (session *session) isToday() bool {
	return session.Date == time.Now().Format(time.DateOnly)
}
