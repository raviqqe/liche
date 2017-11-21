package main

import (
	"os"
	"sync"
)

func main() {
	args, err := getArguments(nil)

	if err != nil {
		fail(err)
	}

	fc := make(chan string, maxOpenFiles)
	ec := make(chan error, 64)
	wg := sync.WaitGroup{}

	go findMarkupFiles(args.filenames, args.recursive, fc, ec)

	wg.Add(1)
	go func() {
		for e := range ec {
			fail(e)
		}

		wg.Done()
	}()

	rc := make(chan fileResult, maxOpenFiles)
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

	wg.Wait()

	if !ok {
		os.Exit(1)
	}
}
