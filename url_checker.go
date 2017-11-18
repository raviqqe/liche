package main

import (
	"net/http"
	"sync"
	"time"
)

type urlChecker struct {
	client    http.Client
	semaphore semaphore
}

func newURLChecker(t time.Duration, s semaphore) urlChecker {
	return urlChecker{http.Client{Timeout: t}, s}
}

func (c urlChecker) Check(s string) (resultErr error) {
	c.semaphore.Request()
	defer c.semaphore.Release()

	res, err := c.client.Get(s)

	if err != nil && res != nil {
		defer func() {
			if err := res.Body.Close(); err != nil {
				resultErr = err
			}
		}()
	}

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
