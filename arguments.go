package main

import "github.com/docopt/docopt-go"

const usage = `Link checker for Markdown and HTML

Usage:
	linkcheck [-v] <filenames>...

Options:
	-v, --verbose  Be verbose`

type arguments struct {
	filenames []string
	verbose   bool
}

func getArgs() (arguments, error) {
	args, err := docopt.Parse(usage, nil, true, "linkcheck", true)

	if err != nil {
		return arguments{}, err
	}

	return arguments{args["<filenames>"].([]string), args["--verbose"].(bool)}, nil
}
