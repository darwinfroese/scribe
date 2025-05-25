package task

import "time"

const (
	sortOrderDesc = iota
	sortOrderAsc
)

func (service *Service) sortOrderCompletedDateFunc(sortOrder int) func(a, b int) int {
	return func(a, b int) int {
		taskA := service.getTask(a).CompletedAt.Round(24 * time.Hour)
		taskB := service.getTask(b).CompletedAt.Round(24 * time.Hour)

		order := taskA.Compare(taskB)

		if sortOrder == sortOrderDesc {
			return order * -1
		}

		return order
	}
}
