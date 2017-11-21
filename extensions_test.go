package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsMarkupFile(t *testing.T) {
	for _, f := range []string{"foo.md", "foo.html", "foo.htm", "foo/bar.md"} {
		assert.True(t, isMarkupFile(f))
	}
}
