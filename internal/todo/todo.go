package todo

import (
	"errors"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

const (
	Command = "todo"

	defaultFolderName = ".todo"
)

// Run is the entrypoint for the todo sub-command
func Run() {
	err := createDirectoryIfNotExists()
	if err != nil {
		log.Fatal(err)
	}
}

func createDirectoryIfNotExists() error {
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	path := filepath.Join(homeDir, defaultFolderName)

	_, err := os.Stat(path)
	if !errors.Is(err, os.ErrNotExist) && err != nil {
		return err
	}

	if err == nil {
		log.Println("using existing '~/.todo' directory")
		return nil
	}

	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		log.Println("created '~/.todo' directory")
	}
	return err
}

// create the ~/.todo/todo main list of todo items
// create the ~/.todo/<date> file
// copy the unfinished TODO items to ~/.todo/<date>
// append the contents of ~/.todo/<date-1> to ~/.todo/todo
