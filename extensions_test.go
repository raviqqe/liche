package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsMarkupFile(t *testing.T) {
	for _, f := range []string{"foo.md", "foo.html", "foo.htm", "foo/bar.md"} {
		assert.True(t, isMarkupFile(f))
	}

	for _, f := range []string{"foo", "foo.m", "bar/foo"} {
		assert.False(t, isMarkupFile(f))
	}
}

func TestIsHTMLFile(t *testing.T) {
	for _, f := range []string{"foo.html", "foo.htm", "foo/bar.html"} {
		assert.True(t, isHTMLFile(f))
	}

	for _, f := range []string{"foo", "foo.md", "bar/foo"} {
		assert.False(t, isHTMLFile(f))
	}
}
