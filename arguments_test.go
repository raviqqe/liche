package main

import (
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetArguments(t *testing.T) {
	for _, c := range []struct {
		argv []string
		args arguments
	}{
		{
			argv: []string{"file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 0, false, nil, false},
		},
		{
			argv: []string{"-c", "42", "file"},
			args: arguments{[]string{"file"}, "", 42, 0, false, nil, false},
		},
		{
			argv: []string{"--concurrency", "42", "file"},
			args: arguments{[]string{"file"}, "", 42, 0, false, nil, false},
		},
		{
			argv: []string{"-d", "directory", "file"},
			args: arguments{[]string{"file"}, "directory", defaultConcurrency, 0, false, nil, false},
		},
		{
			argv: []string{"--document-root", "directory", "file"},
			args: arguments{[]string{"file"}, "directory", defaultConcurrency, 0, false, nil, false},
		},
		{
			argv: []string{"-r", "file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 0, true, nil, false},
		},
		{
			argv: []string{"--recursive", "file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 0, true, nil, false},
		},
		{
			argv: []string{"-t", "42", "file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 42 * time.Second, false, nil, false},
		},
		{
			argv: []string{"--timeout", "42", "file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 42 * time.Second, false, nil, false},
		},
		{
			argv: []string{"-x", "^.*$", "file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 0, false, regexp.MustCompile(`^.*$`), false},
		},
		{
			argv: []string{"--exclude", "^.*$", "file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 0, false, regexp.MustCompile(`^.*$`), false},
		},
		{
			argv: []string{"-v", "file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 0, false, nil, true},
		},
		{
			argv: []string{"--verbose", "file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 0, false, nil, true},
		},
	} {
		args, err := getArguments(c.argv)

		assert.Equal(t, nil, err)
		assert.Equal(t, args, c.args)
	}
}

func TestGetArgumentsWithInvalidArgv(t *testing.T) {
	for _, argv := range [][]string{
		{"-c", "3.14", "file"},
		{"-t", "foo", "file"},
		{"-c", "-t", "file"},
	} {
		_, err := getArguments(argv)
		assert.NotEqual(t, nil, err)
	}
}
