package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/docopt/docopt-go"
)

var defaultConcurrency = func() int {
	const max = 512           // Max number of open files is limited to 1024 on Linux.
	n := 8 * runtime.NumCPU() // 8 is an empirical value.

	if n < max {
		return n
	}

	return max
}()

const usage = `Link checker for Markdown and HTML

Usage:
	liche [-c <num-requests>] [-d <directory>] [-r] [-t <timeout>] [-v] <filenames>...

Options:
	-c, --concurrency <num-requests>  Set max number of concurrent HTTP requests. [default: %v]
	-d, --document-root <directory>  Set document root directory for absolute paths.
	-r, --recursive  Search Markdown and HTML files recursively
	-t, --timeout <timeout>  Set timeout for HTTP requests in seconds. Disabled by default.
	-v, --verbose  Be verbose.`

type arguments struct {
	filenames    []string
	documentRoot string
	concurrency  int
	timeout      time.Duration
	recursive    bool
	verbose      bool
}

func getArguments() (arguments, error) {
	args, err := docopt.Parse(fmt.Sprintf(usage, defaultConcurrency), nil, true, "liche", true)

	if err != nil {
		return arguments{}, err
	}

	c, err := strconv.ParseInt(args["--concurrency"].(string), 10, 32)

	if err != nil {
		return arguments{}, err
	}

	t := 0.0

	if args["--document-root"] == nil {
		args["--document-root"] = ""
	}

	if args["--timeout"] != nil {
		t, err = strconv.ParseFloat(args["--timeout"].(string), 64)

		if err != nil {
			return arguments{}, err
		}
	}

	return arguments{
		args["<filenames>"].([]string),
		args["--document-root"].(string),
		int(c),
		time.Duration(t) * time.Second,
		args["--recursive"].(bool),
		args["--verbose"].(bool),
	}, nil
}
