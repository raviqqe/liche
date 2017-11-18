package main

import (
	"net/url"
	"os"
	"time"

	"github.com/docopt/docopt-go"
)

const usage = `Link checker for Markdown and HTML

Usage:
	linkcheck [-v] <filenames>...

Options:
	-v, --verbose  Be verbose`

func main() {
	defer func() {
		if r := recover(); r != nil {
			printToStderr(r.(error).Error())
			os.Exit(1)
		}
	}()

	args := getArgs()
	fs := args["<filenames>"].([]string)
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
		verbose := args["--verbose"].(bool)

		if !r.Ok() {
			ok = false
			printToStderr(r.String(verbose))
		} else if r.Ok() && verbose {
			printToStderr(r.String(true))
		}
	}

	if !ok {
		os.Exit(1)
	}
}

func isURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https")
}

func getArgs() map[string]interface{} {
	args, err := docopt.Parse(usage, nil, true, "linkcheck", true)

	if err != nil {
		panic(err)
	}

	return args
}
