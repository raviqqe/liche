package main

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestFileCheckerCheck(t *testing.T) {
	c := newFileChecker(0, "", nil, false, false, false, newSemaphore(1024))

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

func TestFileCheckerCheckMany(t *testing.T) {
	c := newFileChecker(0, "", nil, false, false, false, newSemaphore(maxOpenFiles))

	for _, fs := range [][]string{
		{"README.md"},
		{"test/foo.md"},
		{"test/foo.html"},
		{"README.md", "test/foo.md", "test/foo.html"},
	} {
		fc := make(chan string, len(fs))

		for _, f := range fs {
			fc <- f
		}

		close(fc)

		rc := make(chan fileResult, maxOpenFiles)

		c.CheckMany(fc, rc)

		assert.Equal(t, len(fs), len(rc))

		for r := range rc {
			assert.True(t, r.Ok())
		}
	}
}

func TestFileCheckerCheckManyWithInvalidFiles(t *testing.T) {
	c := newFileChecker(0, "", nil, false, false, false, newSemaphore(maxOpenFiles))

	for _, fs := range [][]string{
		{"test/absolute_path.md"},
	} {
		fc := make(chan string, len(fs))

		for _, f := range fs {
			fc <- f
		}

		close(fc)

		rc := make(chan fileResult, maxOpenFiles)

		c.CheckMany(fc, rc)

		assert.Equal(t, len(fs), len(rc))

		ok := true

		for r := range rc {
			ok = ok && r.Ok()
		}

		assert.False(t, ok)
	}
}

func TestFileCheckerExtractURLs(t *testing.T) {
	c := newFileChecker(0, "", nil, false, false, false, newSemaphore(42))

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
