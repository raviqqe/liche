package main

import (
	"net/http"
	"sync"
	"time"
)

type urlChecker struct {
	client http.Client
}

func newURLChecker(timeout time.Duration) urlChecker {
	return urlChecker{http.Client{Timeout: timeout}}
}

func (c urlChecker) Check(s string) error {
	_, err := c.client.Get(s)
	return err
}

func (c urlChecker) CheckMany(ss []string, rc chan<- urlResult) {
	wg := sync.WaitGroup{}

	for _, s := range ss {
		wg.Add(1)
		go func(s string) {
			rc <- urlResult{s, c.Check(s)}
			wg.Done()
		}(s)
	}

	wg.Wait()
	close(rc)
}
