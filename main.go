package main

import "os"

const filesCapacity = 1024

func main() {
	args, err := getArguments(nil)

	if err != nil {
		printToStderr(err.Error())
		os.Exit(1)
	}

	fc := make(chan string, filesCapacity)

	go findMarkupFiles(args.filenames, args.recursive, fc)

	rc := make(chan fileResult, filesCapacity)
	s := newSemaphore(args.concurrency)
	c := newFileChecker(args.timeout, args.documentRoot, s)

	go c.CheckMany(fc, rc)

	ok := true

	for r := range rc {
		if !r.Ok() {
			ok = false
			printToStderr(r.String(args.verbose))
		} else if args.verbose {
			printToStderr(r.String(true))
		}
	}

	if !ok {
		os.Exit(1)
	}
}
