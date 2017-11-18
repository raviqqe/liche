package main

import (
	"os"
	"time"
)

func main() {
	args, err := getArgs()

	if err != nil {
		printToStderr(err.Error())
		os.Exit(1)
	}

	fs := args.filenames
	rc := make(chan fileResult, len(fs))
	c := newFileChecker(5 * time.Second)

	go c.CheckMany(fs, rc)

	ok := true

	for r := range rc {
		if !r.Ok() {
			ok = false
			printToStderr(r.String(args.verbose))
		} else if r.Ok() && args.verbose {
			printToStderr(r.String(true))
		}
	}

	if !ok {
		os.Exit(1)
	}
}
