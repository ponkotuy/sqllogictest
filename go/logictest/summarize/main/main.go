package main

import (
	"fmt"
	"os"

	"github.com/dolthub/sqllogictest/go/logictest"
)

// Summarizes the output of a sqllogictest verify run, reading from stdin or a file argument.
//
// Usage:
//   go run logictest/mysql/main/main.go verify ../../test/evidence/in1.test | go run logictest/summarize/main/main.go
//   go run logictest/summarize/main/main.go result.log
func main() {
	var entries []*logictest.ResultLogEntry
	var err error

	if len(os.Args) > 1 {
		entries, err = logictest.ParseResultFile(os.Args[1])
	} else {
		entries, err = logictest.ParseResults(os.Stdin)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing results: %v\n", err)
		os.Exit(1)
	}

	counts := make(map[logictest.ResultType]int)
	for _, e := range entries {
		counts[e.Result]++
	}

	total := len(entries)
	if total == 0 {
		fmt.Println("No results found.")
		return
	}

	type row struct {
		label string
		rt    logictest.ResultType
	}
	rows := []row{
		{"ok", logictest.Ok},
		{"not ok", logictest.NotOk},
		{"skipped", logictest.Skipped},
		{"timeout", logictest.Timeout},
		{"did not run", logictest.DidNotRun},
	}

	fmt.Printf("%-12s %6s %7s\n", "Result", "Count", "Ratio")
	fmt.Println("-----------------------------")
	for _, r := range rows {
		c := counts[r.rt]
		if c > 0 {
			fmt.Printf("%-12s %6d %6.1f%%\n", r.label, c, float64(c)*100/float64(total))
		}
	}
	fmt.Println("-----------------------------")
	fmt.Printf("%-12s %6d\n", "total", total)
}
