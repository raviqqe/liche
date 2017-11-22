package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestURLCheckerCheck(t *testing.T) {
	c := newURLChecker(0, "", newSemaphore(1024))

	for _, u := range []string{"https://google.com", "README.md"} {
		assert.Equal(t, nil, c.Check(u, "README.md"))
	}

	for _, u := range []string{"https://hey-hey-hi-google.com", "READYOU.md", "://"} {
		assert.NotEqual(t, nil, c.Check(u, "README.md"))
	}
}

func TestURLCheckerCheckWithTimeout(t *testing.T) {
	c := newURLChecker(30*time.Second, "", newSemaphore(1024))

	for _, u := range []string{"https://google.com", "README.md"} {
		assert.Equal(t, nil, c.Check(u, "README.md"))
	}

	for _, u := range []string{"https://hey-hey-hi-google.com", "READYOU.md", "://"} {
		assert.NotEqual(t, nil, c.Check(u, "README.md"))
	}
}

func TestURLCheckerCheckMany(t *testing.T) {
	c := newURLChecker(0, "", newSemaphore(1024))

	for _, us := range [][]string{{}, {"https://google.com", "README.md"}} {
		rc := make(chan urlResult, 1024)
		c.CheckMany(us, "README.md", rc)

		for r := range rc {
			assert.NotEqual(t, "", r.url)
			assert.Equal(t, nil, r.err)
		}
	}
}
func TestURLCheckerResolveURL(t *testing.T) {
	f := newURLChecker(0, "", newSemaphore(1024))

	for _, c := range []struct {
		source, target string
		local          bool
	}{
		{"foo", "foo", true},
		{"https://google.com", "https://google.com", false},
	} {
		u, local, err := f.resolveURL(c.source, "foo.md")

		assert.Equal(t, nil, err)
		assert.Equal(t, c.target, u)
		assert.Equal(t, c.local, local)
	}
}

func TestURLCheckerResolveURLWithAbsolutePath(t *testing.T) {
	f := newURLChecker(0, "", newSemaphore(1024))

	u, _, err := f.resolveURL("/foo", "foo.md")

	assert.NotEqual(t, nil, err)
	assert.Equal(t, "", u)
}

func TestURLCheckerResolveURLWithDocumentRoot(t *testing.T) {
	f := newURLChecker(0, "foo", newSemaphore(1024))

	for _, c := range []struct {
		source, target string
		local          bool
	}{
		{"foo", "foo", true},
		{"https://google.com", "https://google.com", false},
		{"/foo", "foo/foo", true},
	} {
		u, local, err := f.resolveURL(c.source, "foo.md")

		assert.Equal(t, nil, err)
		assert.Equal(t, c.target, u)
		assert.Equal(t, c.local, local)
	}
}
