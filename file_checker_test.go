package main

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestFileCheckerExtractURLs(t *testing.T) {
	c := newFileChecker(0, "", newSemaphore(42))

	for _, x := range []struct {
		html    string
		numURLs int
	}{
		{`<a href="https://google.com">Google</a>`, 1},
		{
			`
			<div>
			<a href="https://google.com">Google</a>
			<a href="https://google.com">Google</a>
			</div>
			`,
			1,
		},
		{
			`
			<div>
			<a href="https://google.com">Google</a>
			<a href="https://yahoo.com">Yahoo!</a>
			</div>
			`,
			2,
		},
	} {
		n, err := html.Parse(strings.NewReader(x.html))

		assert.Equal(t, nil, err)

		us, err := c.extractURLs(n)

		assert.Equal(t, nil, err)
		assert.Equal(t, x.numURLs, len(us))
	}
}

func TestURLParse(t *testing.T) {
	u, err := url.Parse("file-path")

	assert.Equal(t, nil, err)
	assert.Equal(t, "", u.Scheme)
}

func TestIsURL(t *testing.T) {
	for _, s := range []string{"http://google.com", "https://google.com", "file-path"} {
		assert.True(t, isURL(s))
	}

	for _, s := range []string{"ftp://foo.com", "file://file-path", "#foo"} {
		assert.False(t, isURL(s))
	}
}
