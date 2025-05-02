package database

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

const (
	defaultFolderName = ".scribe"

	databaseFileName = "scribe"
)

type Database struct {
	path string
}

func New() *Database {
	path, err := createScribeFolderIfNotExists()
	if err != nil {
		log.Fatal("an error occured creating the scribe folder: ", err)
	}

	file, err := getDatabaseFile(path)
	if err != nil {
		log.Fatal("an error occured creating the scribe database file: ", err)
	}

	return &Database{
		path: fmt.Sprintf("%s/%s", path, file.Name()),
	}
}

func (db *Database) Write(content []byte) error {
	return os.WriteFile(db.path, content, 0644)
}

func (db *Database) Read() ([]byte, error) {
	return os.ReadFile(db.path)
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
