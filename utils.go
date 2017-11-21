package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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

func listDirectory(d string, fc chan<- string) error {
	return filepath.Walk(d, func(f string, i os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		b, err := regexp.MatchString("(^\\.)|(/\\.)", f)

		if err != nil {
			return err
		}

		if !i.IsDir() && !b && isMarkupFile(f) {
			fc <- f
		}

		return nil
	})
}

func findMarkupFiles(fs []string, recursive bool, fc chan<- string) {
	for _, f := range fs {
		i, err := os.Stat(f)

		if err != nil {
			fail(err)
		}

		if i.IsDir() && recursive {
			err := listDirectory(f, fc)

			if err != nil {
				fail(err)
			}

		} else if i.IsDir() {
			fail(fmt.Errorf("%v is not a file", f))
		} else {
			fc <- f
		}
	}

	close(fc)
}
