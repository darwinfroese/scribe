package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/darwinfroese/scribe/cmd"
	"github.com/darwinfroese/scribe/cmd/report"
	"github.com/darwinfroese/scribe/cmd/scribe"
)

func main() {
	args := cmd.Args{}

	scribeCommand := flag.NewFlagSet("scribe", flag.ExitOnError)
	scribeCommand.BoolVar(&args.Global, "global", false, "runs scribe with the global database instead of the local database")

	reportCommand := flag.NewFlagSet("report", flag.ExitOnError)
	reportCommand.IntVar(&args.Last, "last", 0, "the number of sessions to report on, starting with the most recent")
	reportCommand.StringVar(&args.Start, "start", "", "the date to start a report from (YYYY-MM-DD format)")
	reportCommand.StringVar(&args.End, "end", "", "the date to end a report at (YYYY-MM-DD format)")
	reportCommand.BoolVar(&args.All, "all", false, "generate a report for all session dates")
	reportCommand.BoolVar(&args.List, "list", false, "list all session dates")
	reportCommand.BoolVar(&args.Global, "global", false, "runs scribe with the global database instead of the local database")

	if len(os.Args) == 1 {
		scribe.Scribe(args)
		return
	}

	switch os.Args[1] {
	case "report":
		// here we want to parse everything after 'report'
		if err := reportCommand.Parse(os.Args[2:]); err != nil {
			panic(err)
		}
		report.Report(args)
		// I just like a line break between the end of output and the command line after
		// an application exits, this is the easiest way to always apply it.
		fmt.Println()
	default:
		// we don't have a sub-command here
		if err := scribeCommand.Parse(os.Args[1:]); err != nil {
			panic(err)
		}
		scribe.Scribe(args)
	}
}
