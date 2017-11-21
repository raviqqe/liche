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

func listFiles(d string) ([]string, error) {
	fs := []string{}

	err := filepath.Walk(d, func(f string, i os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		b, err := regexp.MatchString("(^\\.)|(/\\.)", f)

		if err != nil {
			return err
		}

		if !i.IsDir() && !b && isMarkupFile(f) {
			fs = append(fs, f)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return fs, nil
}
