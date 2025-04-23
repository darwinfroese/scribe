package notes

import "log"

const (
	Command = "notes"
)

// Run is the entrypoint for the notes sub-command
func Run() {
	log.Println("running notes subcommand")
}
