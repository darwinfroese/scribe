package main

import (
	"log"
	"os"

	"github.com/darwinfroese/scribe/internal/notes"
	"github.com/darwinfroese/scribe/internal/todo"
)

func main() {
	args := os.Args

	if len(args) == 1 {
		log.Fatal("please use one of the commands 'note' or 'todo'")
	}

	switch args[1] {
	case notes.Command:
		notes.Run()
	case todo.Command:
		todo.Run()
	}
}
