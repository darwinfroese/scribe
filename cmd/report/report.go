package report

import (
	"fmt"

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
	fmt.Println("list not implemented")
}

func (svc *service) reportAllSessions() {
	fmt.Println("all not implemented")
}

func (svc *service) reportLastSessions(last int) {
	fmt.Printf("last not implemented (%d)\n", last)
}

func (svc *service) reportDateRangeSessions(start, end string) {
	fmt.Printf("date-range not implemented (%s - %s)\n", start, end)
}
