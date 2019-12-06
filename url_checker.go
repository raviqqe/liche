package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

type urlChecker struct {
	timeout         time.Duration
	documentRoot    string
	excludedPattern *regexp.Regexp
	semaphore       semaphore
}

func newURLChecker(t time.Duration, d string, r *regexp.Regexp, s semaphore) urlChecker {
	return urlChecker{t, d, r, s}
}

func (c urlChecker) Check(u string, f string) error {
	u, local, err := c.resolveURL(u, f)
	if err != nil {
		return err
	}

	if c.excludedPattern != nil && c.excludedPattern.MatchString(u) {
		return nil
	}

	if local {
		_, err := os.Stat(u)
		return err
	}

	c.semaphore.Request()
	defer c.semaphore.Release()

	var sc int
	if c.timeout == 0 {
		sc, _, err = fasthttp.Get(nil, u)
	} else {
		sc, _, err = fasthttp.GetTimeout(nil, u, c.timeout)
	}
	if sc >= http.StatusBadRequest {
		return fmt.Errorf("%s (HTTP error %d)", http.StatusText(sc), sc)
	}
	// Ignore errors from fasthttp about small buffer for URL headers,
	// the content is discarded anyway.
	if _, ok := err.(*fasthttp.ErrSmallBuffer); ok {
		err = nil
	}
	return err
}

func (c urlChecker) CheckMany(us []string, f string, rc chan<- urlResult) {
	wg := sync.WaitGroup{}

	for _, s := range us {
		wg.Add(1)

		go func(s string) {
			rc <- urlResult{s, c.Check(s, f)}
			wg.Done()
		}(s)
	}

	wg.Wait()
	close(rc)
}

func (c urlChecker) resolveURL(u string, f string) (string, bool, error) {
	uu, err := url.Parse(u)

	if err != nil {
		return "", false, err
	}

	if uu.Scheme != "" {
		return u, false, nil
	}

	if !path.IsAbs(uu.Path) {
		return path.Join(filepath.Dir(f), uu.Path), true, nil
	}

	if c.documentRoot == "" {
		return "", false, fmt.Errorf("document root directory is not specified")
	}

	return path.Join(c.documentRoot, uu.Path), true, nil
}
