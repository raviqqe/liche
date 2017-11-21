package main

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestFileCheckerCheck(t *testing.T) {
	c := newFileChecker(0, "", newSemaphore(1024))

	for _, f := range []string{"README.md", "test/foo.md", "test/foo.html"} {
		rs, err := c.Check(f)

		assert.NotEqual(t, 0, len(rs))
		assert.Equal(t, nil, err)

		for _, r := range rs {
			assert.Equal(t, nil, r.err)
		}
	}

	for _, f := range []string{"READYOU.md", "test"} {
		rs, err := c.Check(f)

		assert.Equal(t, ([]urlResult)(nil), rs)
		assert.NotEqual(t, nil, err)
	}

	for _, f := range []string{"test/bad.md", "test/bad.html"} {
		rs, err := c.Check(f)

		assert.Equal(t, nil, err)

		ok := true

		for _, r := range rs {
			if r.err != nil {
				ok = false
			}
		}

		assert.False(t, ok)
	}
}

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

func TestFileCheckerExtractURLsWithInvalidHTML(t *testing.T) {
	c := newFileChecker(0, "", newSemaphore(42))

	for _, s := range []string{
		`<a href="/foo.html">link</a>`,
	} {
		n, err := html.Parse(strings.NewReader(s))

		assert.Equal(t, nil, err)

		us, err := c.extractURLs(n)

		assert.Equal(t, ([]string)(nil), us)
		assert.NotEqual(t, nil, err)
	}
}

func TestFileCheckerResolveURL(t *testing.T) {
	f := newFileChecker(0, "", newSemaphore(1024))

	for _, c := range []struct{ source, target string }{
		{"foo", "foo"},
		{"https://google.com", "https://google.com"},
	} {
		u, err := f.resolveURL(c.source)

		assert.Equal(t, nil, err)
		assert.Equal(t, c.target, u)
	}
}

func TestFileCheckerResolveURLWithAbsolutePath(t *testing.T) {
	f := newFileChecker(0, "", newSemaphore(1024))

	u, err := f.resolveURL("/foo")

	assert.NotEqual(t, nil, err)
	assert.Equal(t, "", u)
}

func TestFileCheckerResolveURLWithDocumentRoot(t *testing.T) {
	f := newFileChecker(0, "foo", newSemaphore(1024))

	for _, c := range []struct{ source, target string }{
		{"foo", "foo"},
		{"https://google.com", "https://google.com"},
		{"/foo", "foo/foo"},
	} {
		u, err := f.resolveURL(c.source)

		assert.Equal(t, nil, err)
		assert.Equal(t, c.target, u)
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
