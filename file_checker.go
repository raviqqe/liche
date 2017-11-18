package main

import (
	"bytes"
	"io/ioutil"
	"time"

	"golang.org/x/net/html"
	"gopkg.in/russross/blackfriday.v2"
)

type fileChecker struct {
	urlChecker urlChecker
}

func newFileChecker(timeout time.Duration, verbose bool) fileChecker {
	return fileChecker{newURLChecker(timeout, verbose)}
}

func (c fileChecker) Check(f string) bool {
	bs, err := ioutil.ReadFile(f)

	if err != nil {
		printToStderr(err.Error())
		return false
	}

	n, err := html.Parse(bytes.NewReader(blackfriday.Run(bs)))

	if err != nil {
		printToStderr(err.Error())
		return false
	}

	return c.urlChecker.CheckMany(extractURLs(n))
}
