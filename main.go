package main

import (
	"fmt"
	"os"

	docopt "github.com/docopt/docopt-go"
)

func main() {
	getArgs()
}

func getArgs() map[string]interface{} {
	usage := `Link checker for Markdown and HTML

Usage:
	linkcheck <filename>`

	args, err := docopt.Parse(usage, nil, true, "linkcheck", true)

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	return args
}
