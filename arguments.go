package main

import (
	"strconv"
	"time"

	"github.com/docopt/docopt-go"
)

const usage = `Link checker for Markdown and HTML

Usage:
	liche [-t <timeout>] [-v] <filenames>...

Options:
	-v, --verbose  Be verbose
	-t, --timeout <timeout>  Set timeout for HTTP requests in seconds [default: 5]`

type arguments struct {
	filenames []string
	timeout   time.Duration
	verbose   bool
}

func getArgs() (arguments, error) {
	args, err := docopt.Parse(usage, nil, true, "liche", true)

	if err != nil {
		return arguments{}, err
	}

	f, err := strconv.ParseFloat(args["--timeout"].(string), 64)

	if err != nil {
		return arguments{}, err
	}

	return arguments{
		args["<filenames>"].([]string),
		time.Duration(f) * time.Second,
		args["--verbose"].(bool),
	}, nil
}
