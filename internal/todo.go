package internal

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

type task struct {
	ID          int    `json:"id"`
	Parent      int    `json:"parent"`
	Priority    int    `json:"priority"`
	Description string `json:"description"`
	Children    []int  `json:"children"`
}

var taskList = []task{}

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
