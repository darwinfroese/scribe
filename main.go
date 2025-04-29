package main

import (
	"errors"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

const (
	defaultFolderName = ".scribe"

	databaseFileName = "scribe"

	priorityLow = iota
	priorityMedium
	priorityHigh
	priorityCritical
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

func main() {
	path, err := createScribeFolderIfNotExists()
	if err != nil {
		log.Fatal("error creating the scribe folder: ", err)
	}

	_, err = getDatabaseFile(path)
	if err != nil {
		log.Fatal("error getting the database file: ", err)
	}
}

func createScribeFolderIfNotExists() (string, error) {
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	path := filepath.Join(homeDir, defaultFolderName)

	_, err := os.Stat(path)
	if !errors.Is(err, os.ErrNotExist) && err != nil {
		return path, err
	}

	if err == nil {
		log.Printf(`using existing scribe directory in "%s"`, path)
		return path, nil
	}

	err = os.Mkdir(path, os.ModePerm)
	if err == nil {
		log.Printf(`created new scribe directory in "%s"`, path)
	}
	return path, err
}

func getDatabaseFile(path string) (os.FileInfo, error) {
	dbPath := filepath.Join(path, databaseFileName)

	file, err := os.Stat(dbPath)
	if errors.Is(err, os.ErrNotExist) {
		err = createDatabaseFile(dbPath)
		if err != nil {
			log.Fatal("error creating the database file: ", err)
		}

		return os.Stat(dbPath)
	}

	return file, err
}

func createDatabaseFile(path string) error {
	_, err := os.Stat(path)
	if !errors.Is(err, os.ErrNotExist) && err != nil {
		return err
	}

	if err == nil {
		log.Printf(`using existing scribe database file at "%s"`, path)
		return nil
	}

	file, err := os.Create(path)
	if err != nil {
		log.Fatal("error creating the database file: ", err)
	}

	err = file.Close() // we aren't using it now so just close it
	if err != nil {
		log.Fatal("error closing the database file: ", err)
	}

	if err == nil {
		log.Printf(`created new scribe database file at "%s"`, path)
	}
	return err
}

// create the ~/.todo/todo main list of todo items
// create the ~/.todo/<date> file
// copy the unfinished TODO items to ~/.todo/<date>
// append the contents of ~/.todo/<date-1> to ~/.todo/todo
