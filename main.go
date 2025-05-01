package main

import (
	"github.com/darwinfroese/scribe/internal/task"
	"github.com/darwinfroese/scribe/internal/ui"
)

func main() {
	// path, err := database.CreateScribeFolderIfNotExists()
	// if err != nil {
	// 	log.Fatal("error creating the scribe folder: ", err)
	// }
	//
	// _, err = database.GetDatabaseFile(path)
	// if err != nil {
	// 	log.Fatal("error getting the database file: ", err)
	// }

	taskService := task.NewTaskService()

	app := ui.New(taskService)
	app.Run()
}

// create the ~/.todo/todo main list of todo items
// create the ~/.todo/<date> file
// copy the unfinished TODO items to ~/.todo/<date>
// append the contents of ~/.todo/<date-1> to ~/.todo/todo
