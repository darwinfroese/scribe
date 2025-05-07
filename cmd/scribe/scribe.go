package main

import (
	"flag"

	"github.com/darwinfroese/scribe/internal/database"
	"github.com/darwinfroese/scribe/internal/task"
	"github.com/darwinfroese/scribe/internal/ui"
)

func main() {
	var global bool
	flag.BoolVar(&global, "global", false, "runs scribe with the global database instead of the local database")
	flag.Parse()

	db := database.New(global)
	taskService := task.NewService(db)

	app := ui.New(taskService)
	app.Run()
}
