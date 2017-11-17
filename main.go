package main

import (
	"fmt"
	"io/ioutil"
	"os"

	docopt "github.com/docopt/docopt-go"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, r.(error).Error())
			os.Exit(1)
		}
	}()

	args := getArgs()

	_, err := ioutil.ReadFile(args["<filename>"].(string))

	if err != nil {
		panic(err)
	}
}

func getArgs() map[string]interface{} {
	usage := `Link checker for Markdown and HTML

Usage:
	linkcheck <filename>`

	args, err := docopt.Parse(usage, nil, true, "linkcheck", true)

	if err != nil {
		panic(err)
	}

	return args
}
