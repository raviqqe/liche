package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/fatih/color"
)

type urlChecker struct {
	client  http.Client
	verbose bool
}

func newURLChecker(timeout time.Duration, verbose bool) urlChecker {
	return urlChecker{http.Client{Timeout: timeout}, verbose}
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

type urlResult struct {
	url string
	err error
}

func (r urlResult) String() string {
	if r.err == nil {
		return color.GreenString("OK") + "\t" + r.url
	}

	return color.RedString("ERROR") + "\t" + r.url + "\t" + color.YellowString(r.err.Error())
}
