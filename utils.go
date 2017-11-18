package main

import (
	"fmt"
	"os"

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

func indent(s string) string {
	return text.Indent(s, "\t")
}
