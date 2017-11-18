package main

import "os"

func main() {
	args, err := getArgs()

	if err != nil {
		printToStderr(err.Error())
		os.Exit(1)
	}

	rc := make(chan fileResult, len(args.filenames))
	c := newFileChecker(args.timeout)

	go c.CheckMany(args.filenames, rc)

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
