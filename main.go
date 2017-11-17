package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/a8m/mark"
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

	bs, err := ioutil.ReadFile(args["<filename>"].(string))

	if err != nil {
		panic(err)
	}

	mark.Render(string(bs))
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
