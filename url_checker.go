package main

import (
	"net/http"
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

func (c urlChecker) Check(s string) bool {
	_, err := c.client.Get(s)

	if s := color.New(color.FgCyan).SprintFunc()(s); err != nil {
		printToStderr(
			color.New(color.FgRed).SprintFunc()("ERROR") + "\t" + s + "\t" + err.Error())
	} else if err == nil && c.verbose {
		printToStderr(color.New(color.FgGreen).SprintFunc()("OK") + "\t" + s)
	}

	return err == nil
}

func (c urlChecker) CheckMany(ss []string) bool {
	bs := make(chan bool, len(ss))

	for _, s := range ss {
		go func(s string) {
			bs <- c.Check(s)
		}(s)
	}

	ok := true

	for i := 0; i < len(ss); i++ {
		ok = ok && <-bs
	}

	return ok
}
