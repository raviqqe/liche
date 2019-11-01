package main

import (
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestURLCheckerCheck(t *testing.T) {
	c := newURLChecker(0, "", nil, false, false, false, newSemaphore(1024))

	for _, u := range []string{"https://google.com", "README.md"} {
		assert.Equal(t, nil, c.Check(u, "README.md"))
	}

	for _, u := range []string{"https://hey-hey-hi-google.com", "READYOU.md", "://"} {
		assert.NotEqual(t, nil, c.Check(u, "README.md"))
	}
}

func TestURLCheckerCheckWithExclude(t *testing.T) {
	c := newURLChecker(0, "", regexp.MustCompile(`^http:\/\/localhost:[13]$`), false, false, false, newSemaphore(1024))

	for _, u := range []string{"http://localhost:1", "http://localhost:3", "README.md"} {
		assert.Equal(t, nil, c.Check(u, "README.md"))
	}

	for _, u := range []string{"http://localhost:2", "READYOU.md"} {
		assert.NotEqual(t, nil, c.Check(u, "README.md"))
	}
}

func TestURLCheckerCheckWithExcludePrivateHosts(t *testing.T) {
	c := newURLChecker(0, "", nil, true, false, false, newSemaphore(1024))

	for _, u := range []string{
		"http://192.168.99.100",
		"http://example.test",
		"http://example.abcdxyz",
	} {
		assert.Equal(t, nil, c.Check(u, "README.md"))
	}
}

func TestURLCheckerCheckWithExcludeLocalhost(t *testing.T) {
	c := newURLChecker(0, "", nil, false, true, false, newSemaphore(1024))

	for _, u := range []string{
		"http://localhost:1",
		"http://localhost:3",
		"http://127.0.0.1:1",
	} {
		assert.Equal(t, nil, c.Check(u, "README.md"))
	}
}

func TestURLCheckerCheckWithExcludeLinkLocal(t *testing.T) {
	c := newURLChecker(0, "", nil, false, false, true, newSemaphore(1024))

	for _, u := range []string{
		"http://169.254.169.254:1",
	} {
		assert.Equal(t, nil, c.Check(u, "README.md"))
	}
}

func TestURLCheckerCheckWithTimeout(t *testing.T) {
	c := newURLChecker(30*time.Second, "", nil, false, false, false, newSemaphore(1024))

	for _, u := range []string{"https://google.com", "README.md"} {
		assert.Equal(t, nil, c.Check(u, "README.md"))
	}

	for _, u := range []string{"https://hey-hey-hi-google.com", "READYOU.md", "://"} {
		assert.NotEqual(t, nil, c.Check(u, "README.md"))
	}
}

func TestURLCheckerCheckMany(t *testing.T) {
	c := newURLChecker(0, "", nil, false, false, false, newSemaphore(1024))

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
	f := newURLChecker(0, "", nil, false, false, false, newSemaphore(1024))

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
	f := newURLChecker(0, "", nil, false, false, false, newSemaphore(1024))

	u, _, err := f.resolveURL("/foo", "foo.md")

	assert.NotEqual(t, nil, err)
	assert.Equal(t, "", u)
}

func TestURLCheckerResolveURLWithDocumentRoot(t *testing.T) {
	f := newURLChecker(0, "foo", nil, false, false, false, newSemaphore(1024))

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
