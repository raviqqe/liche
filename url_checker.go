package main

import (
	"errors"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

type urlChecker struct {
	timeout      time.Duration
	documentRoot string
	semaphore    semaphore
}

func newURLChecker(t time.Duration, d string, s semaphore) urlChecker {
	return urlChecker{t, d, s}
}

func (c urlChecker) Check(u string, f string) error {
	u, err := c.resolveURL(u)

	if err != nil {
		return err
	}

	uu, err := url.Parse(u)

	if err != nil {
		return err
	}

	if uu.Scheme == "" {
		return checkRelativePath(u, f)
	}

	c.semaphore.Request()
	defer c.semaphore.Release()

	if c.timeout == 0 {
		_, _, err := fasthttp.Get(nil, u)
		return err
	}

	_, _, err = fasthttp.GetTimeout(nil, u, c.timeout)
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

func (c urlChecker) resolveURL(u string) (string, error) {
	abs := strings.HasPrefix(u, "/")

	if abs && c.documentRoot != "" {
		return path.Join(c.documentRoot, u), nil
	} else if abs {
		return "", errors.New("document root directory is not specified")
	}

	return u, nil
}

func checkRelativePath(p string, f string) error {
	_, err := os.Stat(path.Join(path.Dir(f), p))
	return err
}
