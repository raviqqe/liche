package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/docopt/docopt-go"
)

var defaultConcurrency = func() int {
	const max = 256
	n := 8 * runtime.NumCPU() // 8 is an empirical value.

	if n < max {
		return n
	}

	return max
}()

const usage = `Link checker for Markdown and HTML

Usage:
	liche [-c <num-requests>] [-t <timeout>] [-v] <filenames>...

Options:
	-c, --concurrency <num-requests>  Set max number of concurrent HTTP requests [default: %v]
	-t, --timeout <timeout>  Set timeout for HTTP requests in seconds [default: 5]
	-v, --verbose  Be verbose`

type arguments struct {
	filenames   []string
	concurrency int
	timeout     time.Duration
	verbose     bool
}

func getArgs() (arguments, error) {
	args, err := docopt.Parse(fmt.Sprintf(usage, defaultConcurrency), nil, true, "liche", true)

	if err != nil {
		return arguments{}, err
	}

	c, err := strconv.ParseInt(args["--concurrency"].(string), 10, 32)

	if err != nil {
		return arguments{}, err
	}

	t, err := strconv.ParseFloat(args["--timeout"].(string), 64)

	if err != nil {
		return arguments{}, err
	}

	return arguments{
		args["<filenames>"].([]string),
		int(c),
		time.Duration(t) * time.Second,
		args["--verbose"].(bool),
	}, nil
}
