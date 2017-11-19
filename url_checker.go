package main

import (
	"net/url"
	"os"
	"path"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

type urlChecker struct {
	timeout   time.Duration
	semaphore semaphore
}

func newURLChecker(t time.Duration, s semaphore) urlChecker {
	return urlChecker{t, s}
}

func (c urlChecker) Check(u string, f string) error {
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

func checkRelativePath(p string, f string) error {
	_, err := os.Stat(path.Join(path.Dir(f), p))
	return err
}
