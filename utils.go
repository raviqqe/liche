package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/kr/text"
)

func stringSetToSlice(s2b map[string]bool) []string {
	ss := make([]string, 0, len(s2b))

	for s := range s2b {
		ss = append(ss, s)
	}

	return ss
}

func printToStderr(xs ...interface{}) {
	fmt.Fprintln(os.Stderr, xs...)
}

func fail(err error) {
	s := err.Error()
	printToStderr(color.RedString(strings.ToUpper(s[:1]) + s[1:]))
	os.Exit(1)
}

func indent(s string) string {
	return text.Indent(s, "\t")
}
