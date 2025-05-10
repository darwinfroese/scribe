package scribe

import (
	"github.com/darwinfroese/scribe/cmd"
	"github.com/darwinfroese/scribe/internal/database"
	"github.com/darwinfroese/scribe/internal/task"
	"github.com/darwinfroese/scribe/internal/ui"
)

func Scribe(args cmd.Args) {
	db := database.New(args.Global)
	taskService := task.NewService(db)

	app := ui.New(taskService)
	app.Run()
}
