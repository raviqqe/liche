package main

import (
	"os"
	"time"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			printToStderr(r.(error).Error())
			os.Exit(1)
		}
	}()

	args, err := getArgs()

	if err != nil {
		panic(err)
	}

	fs := args.filenames
	rc := make(chan fileResult, len(fs))
	c := newFileChecker(5 * time.Second)

	for _, f := range fs {
		go func(f string) {
			rs, err := c.Check(f)

			if err != nil {
				rc <- fileResult{filename: f, err: err}
			}

			rc <- fileResult{filename: f, urlResults: rs}
		}(f)
	}

	ok := true

	for i := 0; i < len(fs); i++ {
		r := <-rc

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
