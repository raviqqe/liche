package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
)

func main() {
	args, err := getArguments(nil)

	if err != nil {
		fail(err)
	}

	m := newMarkupFileFinder()
	wg := sync.WaitGroup{}

	go m.Find(args.filenames, args.recursive)

	wg.Add(1)
	go func() {
		for e := range m.Errors() {
			fail(e)
		}

		wg.Done()
	}()

	rc := make(chan fileResult, maxOpenFiles)
	s := newSemaphore(args.concurrency)
	c := newFileChecker(args.timeout, args.documentRoot, s)

	go c.CheckMany(m.Filenames(), rc)

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

func printToStderr(xs ...interface{}) {
	fmt.Fprintln(os.Stderr, xs...)
}

func fail(err error) {
	s := err.Error()
	printToStderr(color.RedString(strings.ToUpper(s[:1]) + s[1:]))
	os.Exit(1)
}
