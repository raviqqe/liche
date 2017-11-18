package main

import (
	"net/http"
	"runtime"
	"sync"
	"time"
)

const maxOpenFiles = 512

var sem = make(chan bool, func() int {
	n := 8 * runtime.NumCPU() // 8 is an empirical value.

	if n < maxOpenFiles {
		return n
	}

	return maxOpenFiles
}())

type urlChecker struct {
	client http.Client
}

func newURLChecker(timeout time.Duration) urlChecker {
	return urlChecker{http.Client{Timeout: timeout}}
}

func (c urlChecker) Check(s string) error {
	sem <- true
	defer func() { <-sem }()

	res, err := c.client.Get(s)

	if err != nil && res != nil {
		defer res.Body.Close()
	}

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
