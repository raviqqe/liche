package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

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

func listFiles(d string) []string {
	fc := make(chan string, 1024)

	go func() {
		filepath.Walk(d, func(p string, f os.FileInfo, err error) error {
			b, err := regexp.MatchString("(^\\.)|(/\\.)", p)

			if err != nil {
				return err
			}

			if !f.IsDir() && !b && isMarkupFile(p) {
				fc <- p
			}

			return nil
		})

		close(fc)
	}()

	fs := []string{}

	for f := range fc {
		fs = append(fs, f)
	}

	return fs
}
