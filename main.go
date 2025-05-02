package main

import (
	"github.com/darwinfroese/scribe/internal/database"
	"github.com/darwinfroese/scribe/internal/task"
	"github.com/darwinfroese/scribe/internal/ui"
)

func main() {
	db := database.New()
	taskService := task.NewTaskService(db)

	app := ui.New(taskService)
	app.Run()
}
