package task

import (
	"strings"
	"time"
)

const (
	sortOrderDesc = iota
	sortOrderAsc
)

func (service *Service) sortOrderCompletedDateFunc(sortOrder int) func(a, b int) int {
	return func(a, b int) int {
		taskA := service.getTask(a)
		taskB := service.getTask(b)

		taskADate := taskA.CompletedAt.Round(24 * time.Hour)
		taskBDate := taskB.CompletedAt.Round(24 * time.Hour)

		order := taskADate.Compare(taskBDate)

		if order == 0 {
			order = strings.Compare(taskA.Description, taskB.Description)
		}

		if sortOrder == sortOrderDesc {
			return order * -1
		}

		return order
	}
}
