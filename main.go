package main

import (
	"log"

	"github.com/darwinfroese/scribe/internal"
)

func main() {
	path, err := internal.CreateScribeFolderIfNotExists()
	if err != nil {
		log.Fatal("error creating the scribe folder: ", err)
	}

	_, err = internal.GetDatabaseFile(path)
	if err != nil {
		log.Fatal("error getting the database file: ", err)
	}

	internal.Run()
}

// create the ~/.todo/todo main list of todo items
// create the ~/.todo/<date> file
// copy the unfinished TODO items to ~/.todo/<date>
// append the contents of ~/.todo/<date-1> to ~/.todo/todo
