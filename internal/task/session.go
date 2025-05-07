package task

import (
	"slices"
	"time"
)

type session struct {
	ID           int    `json:"id"`
	Date         string `json:"date"`
	PlannedTasks []int  `json:"planned_tasks"`
}

type sessionStorage struct {
	NextID   int       `json:"next_id"`
	Sessions []session `json:"sessions"`
}

func (service *Service) GetSessionIDs() []int {
	ids := []int{}

	for _, session := range service.storage.Sessions.Sessions {
		ids = append(ids, session.ID)
	}

	return ids
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
}

func (service *Service) planTask(taskID int) {
	session := service.getOrCreateSession()

	session.PlannedTasks = append(session.PlannedTasks, taskID)
	service.saveSession(session)
}

func (service *Service) unplanTask(taskID int) {
	session := service.getOrCreateSession()

	for idx, task := range session.PlannedTasks {
		if task == taskID {
			session.PlannedTasks = slices.Delete(session.PlannedTasks, idx, idx+1)
			service.saveSession(session)
			return
		}
	}
}

func (service *Service) saveSession(newSession session) {
	for idx, session := range service.storage.Sessions.Sessions {
		if session.ID == newSession.ID {
			service.storage.Sessions.Sessions[idx] = newSession
			return
		}
	}
}

func (service *Service) getOrCreateSession() session {
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
	service.storage.Sessions.Sessions = append(service.storage.Sessions.Sessions, session)

	return session
}
