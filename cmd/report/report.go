package report

import (
	"fmt"
	"strings"

	"github.com/darwinfroese/scribe/cmd"
	"github.com/darwinfroese/scribe/internal/database"
	"github.com/darwinfroese/scribe/internal/task"
)

type service struct {
	db    *database.Database
	tasks *task.Service
}

func Report(args cmd.Args) {
	db := database.New(args.Global)
	svc := &service{
		db:    db,
		tasks: task.NewService(db),
	}

	if args.List {
		svc.listAllSessions()
		return
	}

	if args.All {
		svc.reportAllSessions()
		return
	}

	if args.Start != "" {
		if args.End == "" {
			svc.reportDateRangeSessions(args.Start, args.Start)
			return
		}

		// TODO: handle args.End being before args.Start
		svc.reportDateRangeSessions(args.Start, args.End)
		return
	}

	if args.Last == 0 {
		args.Last = 1
	}

	svc.reportLastSessions(args.Last)
}

func (svc *service) listAllSessions() {
	dates := svc.tasks.GetAllSessionDates()

	printHeader("All Recorded Session Dates")

	for _, date := range dates {
		fmt.Printf("\t%s\n", date)
	}
}

func (svc *service) reportAllSessions() {
	sessions := svc.tasks.GetAllSessionIDs()

	printHeader("All Sessions")

	for _, session := range sessions {
		svc.printSessionDetails(session)
	}
}

func (svc *service) reportLastSessions(lastCount int) {
	sessions := svc.tasks.GetAllSessionIDs()

	if lastCount > len(sessions) {
		lastCount = len(sessions)
	}

	filtered := sessions[len(sessions)-lastCount:]

	for _, id := range filtered {
		svc.printSessionDetails(id)
	}
}

func (svc *service) reportDateRangeSessions(start, end string) {
	sessions := svc.tasks.GetAllSessionIDs()

	printHeader(fmt.Sprintf("Sessions Between %s And %s", start, end))

	for _, session := range sessions {
		date := svc.tasks.GetSessionDate(session)

		if date >= start && date <= end {
			svc.printSessionDetails(session)
		}
	}
}

func printHeader(header string) {
	length := len(header)
	divider := strings.Repeat("-", length+2)

	fmt.Printf(" %s \n", header)
	fmt.Println(divider)
}

func (svc *service) printSessionDetails(sessionID int) {
	note := svc.tasks.GetNoteForSession(sessionID)
	completedTasks := svc.tasks.GetCompletedTaskIDsForSession(sessionID)
	incompleteTasks := svc.tasks.GetIncompleteTaskIDsForSession(sessionID)

	printHeader(svc.tasks.SessionDisplayString(sessionID))

	fmt.Printf("summary: %s\n", note)
	svc.printTasks("completed tasks:", completedTasks)
	svc.printTasks("incomplete tasks:", incompleteTasks)
}

func (svc *service) printTasks(header string, tasks []int) {
	fmt.Println(header)

	for _, task := range tasks {
		fmt.Printf("- %s\n", svc.tasks.ReportString(task))
	}

	fmt.Println()
}
