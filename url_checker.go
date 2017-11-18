package main

import (
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

func (c urlChecker) Check(u string) error {
	c.semaphore.Request()
	defer c.semaphore.Release()

	_, _, err := fasthttp.GetTimeout(nil, u, c.timeout)
	return err
}

func (c urlChecker) CheckMany(us []string, rc chan<- urlResult) {
	wg := sync.WaitGroup{}

	for _, s := range us {
		wg.Add(1)

		go func(s string) {
			rc <- urlResult{s, c.Check(s)}
			wg.Done()
		}(s)
	}

	wg.Wait()
	close(rc)
}
