package main

import (
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
			args: arguments{[]string{"file"}, "", defaultConcurrency, 0, false, false},
		},
		{
			argv: []string{"-c", "42", "file"},
			args: arguments{[]string{"file"}, "", 42, 0, false, false},
		},
		{
			argv: []string{"--concurrency", "42", "file"},
			args: arguments{[]string{"file"}, "", 42, 0, false, false},
		},
		{
			argv: []string{"-d", "directory", "file"},
			args: arguments{[]string{"file"}, "directory", defaultConcurrency, 0, false, false},
		},
		{
			argv: []string{"--document-root", "directory", "file"},
			args: arguments{[]string{"file"}, "directory", defaultConcurrency, 0, false, false},
		},
		{
			argv: []string{"-r", "file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 0, true, false},
		},
		{
			argv: []string{"--recursive", "file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 0, true, false},
		},
		{
			argv: []string{"-t", "42", "file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 42 * time.Second, false, false},
		},
		{
			argv: []string{"--timeout", "42", "file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 42 * time.Second, false, false},
		},
		{
			argv: []string{"-v", "file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 0, false, true},
		},
		{
			argv: []string{"--verbose", "file"},
			args: arguments{[]string{"file"}, "", defaultConcurrency, 0, false, true},
		},
	} {
		args, err := getArguments(c.argv)

		assert.Equal(t, nil, err)
		assert.Equal(t, args, c.args)
	}
}
